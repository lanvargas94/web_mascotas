export function Line({ className='' }: { className?: string }) {
  return <div className={`animate-pulse bg-gray-200 rounded h-4 ${className}`} />
}

export function CardSkeleton() {
  return (
    <div className="p-3 border rounded bg-white space-y-2">
      <Line className="w-1/3" />
      <Line className="w-2/3" />
    </div>
  )
}

