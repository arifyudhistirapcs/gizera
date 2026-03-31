import db from '@/services/db'
import riskAssessmentService from '@/services/riskAssessmentService'

// Sync status constants
const SYNC_STATUS = {
  SYNCED: 'synced',
  PENDING: 'pending_sync',
  FAILED: 'failed'
}

// ==================== Online/Offline Detection ====================

/**
 * Check if the browser is currently online
 * @returns {boolean}
 */
export function isOnline() {
  return navigator.onLine
}

// ==================== Draft Persistence ====================

/**
 * Save form draft data to IndexedDB.
 * Called on each auto-save from the form view.
 * @param {Object} formData - The full form object (with items)
 */
export async function saveDraftLocally(formData) {
  if (!formData?.id) return

  try {
    const record = {
      formId: formData.id,
      sppgId: formData.sppg_id,
      status: formData.status,
      syncStatus: isOnline() ? SYNC_STATUS.SYNCED : SYNC_STATUS.PENDING,
      updatedAt: new Date().toISOString(),
      data: JSON.stringify(formData)
    }

    const existing = await db.riskAssessmentDrafts
      .where('formId')
      .equals(formData.id)
      .first()

    if (existing) {
      await db.riskAssessmentDrafts.update(existing.id, record)
    } else {
      await db.riskAssessmentDrafts.add(record)
    }
  } catch (err) {
    console.warn('[RiskAssessmentOffline] Failed to save draft locally:', err)
  }
}

/**
 * Load a draft form from IndexedDB by form ID.
 * @param {number} formId
 * @returns {Object|null} The form data or null
 */
export async function loadDraftLocally(formId) {
  try {
    const record = await db.riskAssessmentDrafts
      .where('formId')
      .equals(formId)
      .first()

    if (record?.data) {
      return JSON.parse(record.data)
    }
  } catch (err) {
    console.warn('[RiskAssessmentOffline] Failed to load draft:', err)
  }
  return null
}

// ==================== API Response Cache ====================

/**
 * Cache an API response in IndexedDB.
 * @param {string} cacheKey - Unique key (e.g. 'form_123')
 * @param {Object} data - Response data to cache
 */
export async function cacheResponse(cacheKey, data) {
  try {
    const record = {
      cacheKey,
      data: JSON.stringify(data),
      cachedAt: new Date().toISOString()
    }

    const existing = await db.riskAssessmentCache
      .where('cacheKey')
      .equals(cacheKey)
      .first()

    if (existing) {
      await db.riskAssessmentCache.update(existing.id, record)
    } else {
      await db.riskAssessmentCache.add(record)
    }
  } catch (err) {
    console.warn('[RiskAssessmentOffline] Failed to cache response:', err)
  }
}

/**
 * Get a cached API response from IndexedDB.
 * @param {string} cacheKey
 * @returns {Object|null}
 */
export async function getCachedResponse(cacheKey) {
  try {
    const record = await db.riskAssessmentCache
      .where('cacheKey')
      .equals(cacheKey)
      .first()

    if (record?.data) {
      return JSON.parse(record.data)
    }
  } catch (err) {
    console.warn('[RiskAssessmentOffline] Failed to get cached response:', err)
  }
  return null
}


// ==================== Offline-Aware Save ====================

/**
 * Save draft items — tries API first, falls back to local storage when offline.
 * @param {number} formId
 * @param {Array} items - Array of { id, compliance_score, catatan }
 * @param {Object} fullFormData - The full form object for local caching
 * @returns {{ savedOnline: boolean }}
 */
export async function saveDraft(formId, items, fullFormData) {
  // Always persist locally for resilience
  if (fullFormData) {
    await saveDraftLocally(fullFormData)
  }

  if (isOnline()) {
    try {
      await riskAssessmentService.updateDraft(formId, { items })
      // Mark as synced
      await markDraftSynced(formId)
      return { savedOnline: true }
    } catch (err) {
      console.warn('[RiskAssessmentOffline] API save failed, queued for sync:', err)
      await markDraftPendingSync(formId)
      return { savedOnline: false }
    }
  }

  // Offline — mark for later sync
  await markDraftPendingSync(formId)
  return { savedOnline: false }
}

/**
 * Mark a draft as synced in IndexedDB.
 * @param {number} formId
 */
async function markDraftSynced(formId) {
  try {
    const record = await db.riskAssessmentDrafts
      .where('formId')
      .equals(formId)
      .first()
    if (record) {
      await db.riskAssessmentDrafts.update(record.id, {
        syncStatus: SYNC_STATUS.SYNCED,
        updatedAt: new Date().toISOString()
      })
    }
  } catch (err) {
    console.warn('[RiskAssessmentOffline] Failed to mark draft synced:', err)
  }
}

/**
 * Mark a draft as pending sync in IndexedDB.
 * @param {number} formId
 */
async function markDraftPendingSync(formId) {
  try {
    const record = await db.riskAssessmentDrafts
      .where('formId')
      .equals(formId)
      .first()
    if (record) {
      await db.riskAssessmentDrafts.update(record.id, {
        syncStatus: SYNC_STATUS.PENDING,
        updatedAt: new Date().toISOString()
      })
    }
  } catch (err) {
    console.warn('[RiskAssessmentOffline] Failed to mark draft pending:', err)
  }
}

// ==================== Background Sync ====================

/**
 * Sync all pending_sync drafts to the API.
 * Called when the app comes back online.
 * Conflict resolution: server timestamp wins — if the server already has
 * a newer version, we skip the local update.
 * @returns {{ synced: number, skipped: number, failed: number }}
 */
export async function syncPendingDrafts() {
  const result = { synced: 0, skipped: 0, failed: 0 }

  try {
    const pendingDrafts = await db.riskAssessmentDrafts
      .where('syncStatus')
      .equals(SYNC_STATUS.PENDING)
      .toArray()

    for (const draft of pendingDrafts) {
      try {
        const formData = JSON.parse(draft.data)

        // Skip submitted forms — nothing to sync
        if (formData.status === 'submitted') {
          await markDraftSynced(draft.formId)
          result.skipped++
          continue
        }

        // Fetch server version to check for conflicts
        let serverForm = null
        try {
          const resp = await riskAssessmentService.getForm(draft.formId)
          serverForm = resp.data?.data
        } catch (fetchErr) {
          // 404 means form was deleted on server — skip
          if (fetchErr.response?.status === 404) {
            await db.riskAssessmentDrafts.delete(draft.id)
            result.skipped++
            continue
          }
          throw fetchErr
        }

        // Conflict resolution: server timestamp wins
        if (serverForm) {
          const serverUpdated = new Date(serverForm.updated_at).getTime()
          const localUpdated = new Date(draft.updatedAt).getTime()

          if (serverUpdated > localUpdated) {
            // Server is newer — skip local changes, update local cache with server data
            await saveDraftLocally(serverForm)
            await markDraftSynced(draft.formId)
            result.skipped++
            continue
          }

          // Server form already submitted — skip
          if (serverForm.status === 'submitted') {
            await saveDraftLocally(serverForm)
            await markDraftSynced(draft.formId)
            result.skipped++
            continue
          }
        }

        // Push local changes to API
        const items = (formData.items || []).map(i => ({
          item_id: i.id,
          compliance_score: i.compliance_score,
          catatan: i.catatan || ''
        }))

        await riskAssessmentService.updateDraft(draft.formId, { items })
        await markDraftSynced(draft.formId)
        result.synced++
      } catch (itemErr) {
        console.error(`[RiskAssessmentOffline] Failed to sync form ${draft.formId}:`, itemErr)
        // Mark as failed but don't delete — will retry next time
        try {
          await db.riskAssessmentDrafts.update(draft.id, {
            syncStatus: SYNC_STATUS.FAILED
          })
        } catch (_) { /* ignore */ }
        result.failed++
      }
    }
  } catch (err) {
    console.error('[RiskAssessmentOffline] syncPendingDrafts error:', err)
  }

  return result
}

/**
 * Get count of drafts pending sync.
 * @returns {number}
 */
export async function getPendingSyncCount() {
  try {
    return await db.riskAssessmentDrafts
      .where('syncStatus')
      .equals(SYNC_STATUS.PENDING)
      .count()
  } catch (err) {
    return 0
  }
}

// ==================== Online/Offline Listener ====================

let _onlineHandler = null

/**
 * Setup listener that triggers sync when the app comes back online.
 * Call once at app/view mount.
 * @param {Function} [onSyncComplete] - Optional callback with sync result
 * @returns {Function} cleanup function to remove the listener
 */
export function setupOfflineSync(onSyncComplete) {
  const handler = async () => {
    console.log('[RiskAssessmentOffline] Back online — syncing pending drafts')
    const result = await syncPendingDrafts()
    if (typeof onSyncComplete === 'function') {
      onSyncComplete(result)
    }
  }

  window.addEventListener('online', handler)
  _onlineHandler = handler

  return () => {
    window.removeEventListener('online', handler)
    _onlineHandler = null
  }
}

// ==================== Cleanup ====================

/**
 * Remove cached data older than the given number of days.
 * @param {number} daysToKeep - Default 30
 */
export async function cleanupOldData(daysToKeep = 30) {
  try {
    const cutoff = new Date()
    cutoff.setDate(cutoff.getDate() - daysToKeep)
    const cutoffISO = cutoff.toISOString()

    // Clean old cache entries
    const oldCache = await db.riskAssessmentCache
      .filter(r => r.cachedAt < cutoffISO)
      .toArray()
    if (oldCache.length) {
      await db.riskAssessmentCache.bulkDelete(oldCache.map(r => r.id))
    }

    // Clean old synced drafts (keep pending ones)
    const oldDrafts = await db.riskAssessmentDrafts
      .filter(r => r.syncStatus === SYNC_STATUS.SYNCED && r.updatedAt < cutoffISO)
      .toArray()
    if (oldDrafts.length) {
      await db.riskAssessmentDrafts.bulkDelete(oldDrafts.map(r => r.id))
    }
  } catch (err) {
    console.warn('[RiskAssessmentOffline] Cleanup error:', err)
  }
}
