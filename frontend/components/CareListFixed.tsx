"use client"
import React, { useState } from 'react'
import { useToast } from './Toast'
import { usePetCares, addCare, removeCare } from '../lib/hooks'
import type { Cuidado } from '../lib/api/client'
import { CardSkeleton } from './Skeleton'
import Tooltip from './Tooltip'

export default function CareList({ mascotaId }: { mascotaId: number }) {
  const { cares, isLoading } = usePetCares(mascotaId)
  const defaultTipo = 'Vacunacion' as Cuidado['tipo_cuidado']
  const [form, setForm] = useState<{ tipo_cuidado: Cuidado['tipo_cuidado']; descripcion: string; fecha: string; hora: string }>({ tipo_cuidado: defaultTipo, descripcion: '', fecha: '', hora: '' })
  const { show } = useToast()

  const fmt2 = (n:number) => String(n).padStart(2,'0')
  const todayLocal = () => {
    const d = new Date()
    return `${d.getFullYear()}-${fmt2(d.getMonth()+1)}-${fmt2(d.getDate())}`
  }
  const tomorrowLocal = () => {
    const d = new Date()
    d.setDate(d.getDate()+1)
    return `${d.getFullYear()}-${fmt2(d.getMonth()+1)}-${fmt2(d.getDate())}`
  }
  const isPastOrSameDate = (dateStr:string) => dateStr !== '' && dateStr <= todayLocal()
  const isSunday = (dateStr:string) => {
    if (!dateStr) return false
    const d = new Date(`${dateStr}T00:00:00`)
    return d.getDay() === 0
  }

  const add = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const { fecha, hora, descripcion, tipo_cuidado } = form
      if (isPastOrSameDate(fecha)) { show('Solo puede programar cuidados a partir del dia siguiente.','error'); return }
      if (isSunday(fecha)) { show('No es posible registrar cuidados los dias domingo. Seleccione entre lunes y sabado.','error'); return }
      const iso = new Date(`${fecha}T${hora}`).toISOString()
      await addCare(mascotaId, { tipo_cuidado, descripcion, fecha_cuidado: iso })
      setForm({ tipo_cuidado: defaultTipo, descripcion: '', fecha: '', hora: '' })
      show('Cuidado agregado','success')
    } catch(e:any) { show(e.message || 'Error al agregar','error') }
  }

  return (
    <div className="space-y-3">
      <form onSubmit={add} className="p-3 border rounded bg-white space-y-2">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-2 items-end">
          <select className="border p-2 rounded" value={form.tipo_cuidado} onChange={e=>setForm({...form, tipo_cuidado: e.target.value as Cuidado['tipo_cuidado']})}>
            <option>Vacunacion</option>
            <option>Desparasitacion</option>
            <option>Consulta Veterinaria</option>
            <option>Bano</option>
          </select>
          <div>
            <label className="block text-sm mb-1">Fecha <Tooltip label="Obligatorio. Selecciona la fecha del cuidado." /></label>
            <input type="date" className="border p-2 rounded w-full" min={tomorrowLocal()} value={form.fecha} onChange={e=>setForm({...form, fecha:e.target.value})} />
          </div>
          <div>
            <label className="block text-sm mb-1">Hora <Tooltip label="Obligatorio. Selecciona la hora y minuto." /></label>
            <input type="time" className="border p-2 rounded w-full" value={form.hora} onChange={e=>setForm({...form, hora:e.target.value})} />
          </div>
          <textarea className="border p-2 rounded md:col-span-4" placeholder="Descripcion" value={form.descripcion} onChange={e=>setForm({...form, descripcion:e.target.value})} />
        </div>
        <button disabled={!form.fecha || !form.hora || (form.descripcion.trim().length<2) || isPastOrSameDate(form.fecha) || isSunday(form.fecha)} className="px-3 py-2 bg-green-600 text-white rounded disabled:opacity-50 disabled:cursor-not-allowed">Agregar cuidado</button>
      </form>
      <ul className="space-y-2">
        {isLoading && (<>
          <CardSkeleton /><CardSkeleton />
        </>)}
        {!isLoading && cares.map(i => {
          const d = new Date(i.fecha_cuidado)
          const yyyy = d.getFullYear()
          const mm = String(d.getMonth()+1).padStart(2,'0')
          const dd = String(d.getDate()).padStart(2,'0')
          const dateOnly = `${yyyy}-${mm}-${dd}`
          const hrs = d.getHours()
          const mins = d.getMinutes()
          const hasTime = !isNaN(hrs) && !isNaN(mins) && (hrs !== 0 || mins !== 0)
          const time12 = new Date(d).toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit', hour12: true })
          const when = hasTime ? `${dateOnly}, hora ${time12}` : dateOnly
          return (
          <li key={i.id} className="p-3 border rounded bg-white flex items-center justify-between">
            <div>
              <div className="font-medium">{i.tipo_cuidado}</div>
              <div className="text-sm text-gray-600">{when}</div>
              <div className="text-sm">{i.descripcion}</div>
            </div>
            <button className="text-red-600" onClick={async () => {
              try { await removeCare(mascotaId, i.id); show('Cuidado eliminado','success') }
              catch(e:any){ show(e.message || 'Error al eliminar','error') }
            }}>Eliminar</button>
          </li>
        )})}
      </ul>
    </div>
  )
}

