import '../styles/globals.css'
import React from 'react'
import { ToastProvider } from '../components/Toast'

export const metadata = {
  title: 'Mascotas',
  description: 'Gestión de mascotas y cuidados',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="es">
      <body className="min-h-screen bg-gray-50 text-gray-900">
        <ToastProvider>
          <header className="w-full border-b bg-white">
            <div className="max-w-5xl mx-auto px-4 py-3 sm:px-6">
              <h1 className="text-xl sm:text-2xl font-semibold">Mascotas</h1>
            </div>
          </header>
          <main className="w-full">
            <div className="max-w-5xl mx-auto px-4 py-6 sm:px-6">
              {/* Fallback simple para asegurar contenido mínimo */}
              <div className="sr-only">Aplicación Mascotas</div>
              {children}
            </div>
          </main>
        </ToastProvider>
      </body>
    </html>
  )
}
