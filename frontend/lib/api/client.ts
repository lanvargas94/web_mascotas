export type Mascota = {
  id: number
  nombre: string
  especie: 'Perro' | 'Gato' | 'Conejo'
  raza: string
  fecha_nacimiento: string
  sexo: 'Macho' | 'Hembra'
}

export type Cuidado = {
  id: number
  tipo_cuidado: 'Vacunación' | 'Desparasitación' | 'Consulta Veterinaria' | 'Baño'
  descripcion: string
  fecha_cuidado: string
  mascota_id: number
}

type ApiError = { error: { code: string; message: string; fields?: { field: string; message: string }[] } }

const API = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API}${path}`, { ...init, headers: { 'Content-Type': 'application/json', ...(init?.headers || {}) } })
  if (!res.ok) {
    let err: ApiError | undefined
    try { err = await res.json() } catch {}
    const msg = err?.error?.message || `HTTP ${res.status}`
    throw new Error(msg)
  }
  return res.json() as Promise<T>
}

// Pets API
export const listPets = (limit=50, offset=0) => apiFetch<Mascota[]>(`/mascotas?limit=${limit}&offset=${offset}`)
export const createPet = (data: Omit<Mascota,'id'>) => apiFetch<Mascota>(`/mascotas`, { method: 'POST', body: JSON.stringify(data) })
export const deletePet = (id: number) => fetch(`${API}/mascotas/${id}`, { method: 'DELETE' }).then(r => { if(!r.ok) throw new Error(`HTTP ${r.status}`) })

// Cares API
export const listCares = (petId: number) => apiFetch<Cuidado[]>(`/mascotas/${petId}/cuidados`)
export const createCare = (petId: number, data: Omit<Cuidado,'id'|'mascota_id'>) => apiFetch<Cuidado>(`/mascotas/${petId}/cuidados`, { method: 'POST', body: JSON.stringify(data) })
export const deleteCare = (id: number) => fetch(`${API}/cuidados/${id}`, { method: 'DELETE' }).then(r => { if(!r.ok) throw new Error(`HTTP ${r.status}`) })

