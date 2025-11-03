"use client"
import React from 'react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useToast } from './Toast'
import { addPet } from '../lib/hooks'
import Tooltip from './Tooltip'

const API = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

const schema = z.object({
  nombre: z.string().min(2, 'Minimo 2 caracteres').max(100),
  especie: z.enum(["Perro","Gato","Conejo"]),
  raza: z.string().min(2, 'Minimo 2 caracteres').max(100),
  fecha_nacimiento: z.string()
    .regex(/^\d{4}-\d{2}-\d{2}$/,'Formato YYYY-MM-DD')
    .refine((v)=> {
      const d = new Date(v+ 'T00:00:00')
      const today = new Date()
      // normaliza a medianoche local
      today.setHours(0,0,0,0)
      return d.getTime() <= today.getTime()
    }, 'No puede ser una fecha futura'),
  sexo: z.enum(["Macho","Hembra"]),
})

type FormValues = z.infer<typeof schema>

export default function PetForm() {
  const { show } = useToast()
  const { register, handleSubmit, reset, formState: { errors, isSubmitting } } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { nombre: '', especie: 'Perro', raza: '', fecha_nacimiento: '', sexo: 'Macho' },
  })

  const onSubmit = async (data: FormValues) => {
    try {
      await addPet(data)
      reset()
      show('Mascota creada', 'success')
    } catch (e:any) {
      console.error('Error creando mascota', e)
      show('No se pudo crear la mascota', 'error')
    }
  }

  const today = new Date(); today.setHours(0,0,0,0)
  const todayStr = today.toISOString().slice(0,10)
  return (
    <form onSubmit={handleSubmit(onSubmit)} className="p-4 bg-white border border-gray-200 rounded space-y-3" aria-labelledby="form-title">
      <h2 id="form-title" className="text-xl font-medium">Nueva Mascota</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div>
          <label className="block text-sm mb-1" htmlFor="nombre">Nombre <span className="text-red-600">*</span> <Tooltip label="Campo obligatorio. Minimo 2 caracteres." /></label>
          <input id="nombre" aria-invalid={!!errors.nombre} className="w-full border p-2 rounded" placeholder="Nombre" {...register('nombre')} />
          {errors.nombre && <p className="text-sm text-red-600">{errors.nombre.message as string}</p>}
        </div>
        <div>
          <label className="block text-sm mb-1" htmlFor="raza">Raza <span className="text-red-600">*</span> <Tooltip label="Campo obligatorio. Minimo 2 caracteres." /></label>
          <input id="raza" aria-invalid={!!errors.raza} className="w-full border p-2 rounded" placeholder="Raza" {...register('raza')} />
          {errors.raza && <p className="text-sm text-red-600">{errors.raza.message as string}</p>}
        </div>
        <div>
          <label className="block text-sm mb-1" htmlFor="especie">Especie</label>
          <select id="especie" className="w-full border p-2 rounded" {...register('especie')}>
            <option>Perro</option>
            <option>Gato</option>
            <option>Conejo</option>
          </select>
        </div>
        <div>
          <label className="block text-sm mb-1" htmlFor="sexo">Sexo</label>
          <select id="sexo" className="w-full border p-2 rounded" {...register('sexo')}>
            <option>Macho</option>
            <option>Hembra</option>
          </select>
        </div>
        <div className="sm:col-span-2">
          <label className="block text-sm mb-1" htmlFor="fecha">Fecha de nacimiento <span className="text-red-600">*</span> <Tooltip label="Obligatorio. Formato YYYY-MM-DD. No puede ser futura." /></label>
          <input id="fecha" type="date" max={todayStr} aria-invalid={!!errors.fecha_nacimiento} className="w-full border p-2 rounded" {...register('fecha_nacimiento')} />
          <p className="text-xs text-gray-500 mt-1">Formato YYYY-MM-DD</p>
          {errors.fecha_nacimiento && <p className="text-sm text-red-600">{errors.fecha_nacimiento.message as string}</p>}
        </div>
      </div>
      <button disabled={isSubmitting} className="px-4 py-2 bg-blue-600 text-white rounded">{isSubmitting ? 'Guardando...' : 'Guardar'}</button>
    </form>
  )
}


