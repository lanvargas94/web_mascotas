"use client"
import React, { useState } from 'react'
import CareList from './CareListFixed'
import { useToast } from './Toast'
import { usePets, removePet } from '../lib/hooks'
import { CardSkeleton } from './Skeleton'

type Mascota = {
  id: number
  nombre: string
  especie: string
  raza: string
  fecha_nacimiento: string
  sexo: string
}

const API = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export default function PetList() {
  const { pets, isLoading } = usePets()
  const [selected, setSelected] = useState<Mascota | null>(null)
  const { show } = useToast()

  const fmtDateOnly = (s: string) => {
    // Expect ISO or YYYY-MM-DD; return only YYYY-MM-DD
    if (!s) return ''
    const tIdx = s.indexOf('T')
    return tIdx > 0 ? s.slice(0, tIdx) : s
  }

  return (
    <div>
      <h2 className="text-xl font-medium mb-2">Lista de Mascotas</h2>
      <ul className="space-y-2">
        {isLoading && (<>
          <CardSkeleton /><CardSkeleton /><CardSkeleton />
        </>)}
        {!isLoading && pets.map((p) => (
          <li key={p.id} className={`p-3 border rounded bg-white flex items-center justify-between ${selected?.id===p.id?'border-blue-400':'border-gray-200'}`}>
            <div>
              <div className="font-medium">{p.nombre} <span className="text-sm text-gray-500">({p.especie})</span></div>
              <div className="text-sm text-gray-600">{p.raza} • {p.sexo} • {fmtDateOnly(p.fecha_nacimiento)}</div>
            </div>
            <div className="space-x-2">
              <button className="text-blue-600" onClick={() => setSelected(selected?.id===p.id ? null : p)}>{selected?.id===p.id ? 'Ocultar' : 'Cuidados'}</button>
              <button aria-label={`Eliminar ${p.nombre}`} className="text-red-600" onClick={async () => {
                try { await removePet(p.id); setSelected(null); show('Mascota eliminada','success') }
                catch (e:any) { show(e.message || 'Error al eliminar','error') }
              }}>Eliminar</button>
            </div>
          </li>
        ))}
      </ul>
      {selected && (
        <div className="mt-6">
          <h3 className="text-lg font-medium mb-2">Cuidados de {selected.nombre}</h3>
          <CareList mascotaId={selected.id} />
        </div>
      )}
    </div>
  )
}

