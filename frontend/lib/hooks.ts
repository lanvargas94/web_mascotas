"use client"
import useSWR, { mutate as globalMutate } from 'swr'
import { listPets, listCares, createPet, deletePet, createCare, deleteCare, Mascota, Cuidado } from './api/client'

export function usePets(limit=50, offset=0) {
  const key = ['pets', limit, offset]
  const { data, error, isLoading, mutate } = useSWR(key, () => listPets(limit, offset))
  return { pets: data || [], error, isLoading, mutate }
}

export function usePetCares(petId: number | null) {
  const shouldFetch = !!petId
  const { data, error, isLoading, mutate } = useSWR(shouldFetch ? ['cares', petId] : null, () => listCares(petId as number))
  return { cares: data || [], error, isLoading, mutate }
}

// Mutations
export async function addPet(data: Omit<Mascota,'id'>) {
  const m = await createPet(data)
  await globalMutate((key:any) => Array.isArray(key) && key[0]==='pets')
  return m
}

export async function removePet(id: number) {
  await deletePet(id)
  await globalMutate((key:any) => Array.isArray(key) && key[0]==='pets')
}

export async function addCare(petId: number, data: Omit<Cuidado,'id'|'mascota_id'>) {
  const c = await createCare(petId, data)
  await globalMutate(['cares', petId])
  return c
}

export async function removeCare(petId: number, id: number) {
  await deleteCare(id)
  await globalMutate(['cares', petId])
}

