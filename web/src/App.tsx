import { useEffect, useMemo, useState } from 'react'

type ChecklistItem = {
  id: string
  title: string
  category: string
  done: boolean
}

type ChecklistResponse = {
  items: ChecklistItem[]
}

export function App() {
  const [items, setItems] = useState<ChecklistItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function loadChecklist() {
      try {
        const response = await fetch('/api/checklist')
        if (!response.ok) {
          throw new Error(`Request failed with ${response.status}`)
        }
        const data = (await response.json()) as ChecklistResponse
        setItems(data.items)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unable to load checklist')
      } finally {
        setLoading(false)
      }
    }

    void loadChecklist()
  }, [])

  const completedCount = useMemo(
    () => items.filter((item) => item.done).length,
    [items],
  )

  function toggleItem(id: string) {
    setItems((current) =>
      current.map((item) =>
        item.id === id ? { ...item, done: !item.done } : item,
      ),
    )
  }

  return (
    <main className="container app-shell">
      <header className="summary">
        <div>
          <p className="eyebrow">Track day prep</p>
          <h1>Raceday Checklist</h1>
        </div>
        <div className="progress" aria-label="Checklist progress">
          <progress value={completedCount} max={items.length || 1} />
          <strong>{completedCount}/{items.length}</strong>
          <span>ready</span>
        </div>
      </header>

      {loading && <article aria-busy="true">Loading checklist...</article>}
      {error && <article className="error">API error: {error}</article>}

      {!loading && !error && (
        <section className="checklist" aria-label="Checklist items">
          {items.map((item) => (
            <label className="checklist-item" key={item.id}>
              <input
                type="checkbox"
                checked={item.done}
                onChange={() => toggleItem(item.id)}
              />
              <span>
                <strong>{item.title}</strong>
                <small>{item.category}</small>
              </span>
            </label>
          ))}
        </section>
      )}
    </main>
  )
}
