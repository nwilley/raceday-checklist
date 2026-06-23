import { useEffect, useMemo, useState } from 'react'

type ChecklistItem = {
  id: string
  sectionId: string
  itemId: string
  title: string
  category: string
  done: boolean
}

type ChecklistResponse = {
  items: ChecklistItem[]
}

type ChecklistSection = {
  name: string
  items: Array<ChecklistItem & { index: number }>
}

const clientStorageKey = 'racedayClientId'
const clientHeader = 'X-Raceday-Client'

function getClientId() {
  const existing = window.localStorage.getItem(clientStorageKey)
  if (existing) {
    return existing
  }

  const clientId = window.crypto.randomUUID()
  window.localStorage.setItem(clientStorageKey, clientId)
  return clientId
}

function itemKey(item: Pick<ChecklistItem, 'sectionId' | 'itemId'>) {
  return `${item.sectionId}/${item.itemId}`
}

export function App() {
  const [items, setItems] = useState<ChecklistItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [pendingItems, setPendingItems] = useState<Set<string>>(() => new Set())
  const [saveError, setSaveError] = useState<string | null>(null)

  useEffect(() => {
    async function loadChecklist() {
      try {
        const response = await fetch('/api/checklist', {
          headers: {
            [clientHeader]: getClientId(),
          },
        })
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

  const sections = useMemo(() => {
    const grouped = new Map<string, ChecklistSection>()

    items.forEach((item, index) => {
      const sectionName = item.category || 'Checklist'
      const section = grouped.get(sectionName) ?? {
        name: sectionName,
        items: [],
      }

      section.items.push({ ...item, index })
      grouped.set(sectionName, section)
    })

    return Array.from(grouped.values())
  }, [items])

  async function toggleItem(index: number) {
    const item = items[index]
    if (!item) {
      return
    }

    const key = itemKey(item)
    const previousDone = item.done
    const nextDone = !previousDone

    setSaveError(null)
    setPendingItems((current) => new Set(current).add(key))
    setItems((current) =>
      current.map((currentItem, itemIndex) =>
        itemIndex === index ? { ...currentItem, done: nextDone } : currentItem,
      ),
    )

    try {
      const response = await fetch(
        `/api/checklist/items/${encodeURIComponent(item.sectionId)}/${encodeURIComponent(item.itemId)}`,
        {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
            [clientHeader]: getClientId(),
          },
          body: JSON.stringify({ done: nextDone }),
        },
      )
      if (!response.ok) {
        throw new Error(`Request failed with ${response.status}`)
      }
    } catch (err) {
      setItems((current) =>
        current.map((currentItem, itemIndex) =>
          itemIndex === index
            ? { ...currentItem, done: previousDone }
            : currentItem,
        ),
      )
      setSaveError(
        err instanceof Error ? err.message : 'Unable to save checklist item',
      )
    } finally {
      setPendingItems((current) => {
        const next = new Set(current)
        next.delete(key)
        return next
      })
    }
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
      {saveError && (
        <article className="error">Save failed: {saveError}</article>
      )}

      {!loading && !error && (
        <section className="checklist" aria-label="Checklist sections">
          {sections.map((section) => (
            <section className="checklist-section" key={section.name}>
              <header>
                <h2>{section.name}</h2>
                <span>
                  {section.items.filter((item) => item.done).length}/
                  {section.items.length}
                </span>
              </header>

              <div className="section-items">
                {section.items.map((item) => (
                  <label
                    className="checklist-item"
                    key={`${section.name}-${item.sectionId}-${item.itemId}-${item.index}`}
                  >
                    <input
                      type="checkbox"
                      checked={item.done}
                      disabled={pendingItems.has(itemKey(item))}
                      onChange={() => toggleItem(item.index)}
                    />
                    <span>
                      <strong>{item.title}</strong>
                      {pendingItems.has(itemKey(item)) && (
                        <small>Saving...</small>
                      )}
                    </span>
                  </label>
                ))}
              </div>
            </section>
          ))}
        </section>
      )}
    </main>
  )
}
