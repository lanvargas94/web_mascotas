"use client"
import React from 'react'
import { useToast } from '../components/Toast'

export default function GlobalError({ error, reset }: { error: Error & { digest?: string }, reset: () => void }) {
  const { show } = useToast()
  React.useEffect(() => {
    if (error?.message) {
      console.error('App error:', error)
      show('Ocurrió un error en la aplicación', 'error')
    }
  }, [error, show])
  return (
    <html lang="es">
      <body>
        <div className="min-h-screen flex items-center justify-center p-6">
          <div className="max-w-md text-center">
            <h2 className="text-xl font-semibold mb-2">Algo salió mal</h2>
            <p className="text-sm text-gray-600 mb-4">Intenta recargar la página. Si persiste, avísanos.</p>
            <button onClick={() => reset()} className="px-4 py-2 bg-blue-600 text-white rounded">Reintentar</button>
          </div>
        </div>
      </body>
    </html>
  )
}

