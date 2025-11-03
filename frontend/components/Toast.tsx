"use client"
import React, { createContext, useCallback, useContext, useMemo, useState } from 'react'

type Toast = { id: number; message: string; type?: 'success' | 'error' }

type ToastCtx = {
  show: (message: string, type?: 'success'|'error') => void
}

const Ctx = createContext<ToastCtx | null>(null)

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([])
  const show = useCallback((message: string, type: 'success'|'error' = 'success') => {
    const id = Date.now() + Math.random()
    setToasts((t) => [...t, { id, message, type }])
    if (typeof window !== 'undefined') {
      window.setTimeout(() => setToasts((t) => t.filter(x => x.id !== id)), 3000)
    }
  }, [])
  const value = useMemo(() => ({ show }), [show])
  return (
    <Ctx.Provider value={value}>
      {children}
      <div aria-live="polite" aria-atomic="true" className="fixed bottom-4 right-4 flex flex-col gap-2 z-50">
        {toasts.map(t => (
          <div key={t.id} role="status" className={`px-4 py-2 rounded shadow text-white ${t.type==='error'?'bg-red-600':'bg-green-600'}`}>{t.message}</div>
        ))}
      </div>
    </Ctx.Provider>
  )
}

export function useToast() {
  const ctx = useContext(Ctx)
  if (!ctx) throw new Error('useToast must be used within ToastProvider')
  return ctx
}
