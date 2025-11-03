"use client"
import React from 'react'

export default function Tooltip({ label, children }: { label: string, children?: React.ReactNode }) {
  return (
    <span className="relative inline-flex items-center group">
      {children ?? (
        <span aria-hidden className="ml-1 inline-flex items-center justify-center w-4 h-4 rounded-full bg-gray-300 text-[10px] text-gray-800">i</span>
      )}
      <span role="tooltip" className="pointer-events-none absolute bottom-full mb-1 hidden w-max max-w-xs rounded bg-black px-2 py-1 text-xs text-white opacity-0 group-hover:block group-hover:opacity-100">
        {label}
      </span>
    </span>
  )
}

