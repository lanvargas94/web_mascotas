import dynamic from 'next/dynamic'

const PetForm = dynamic(() => import('../components/PetForm'), { ssr: false, loading: () => <div>Cargando formulario…</div> })
const PetList = dynamic(() => import('../components/PetList'), { ssr: false, loading: () => <div>Cargando listado…</div> })

export default function Page() {
  return (
    <div className="grid gap-8 md:grid-cols-2">
      <div>
        <PetForm />
      </div>
      <div>
        <PetList />
      </div>
    </div>
  )
}
