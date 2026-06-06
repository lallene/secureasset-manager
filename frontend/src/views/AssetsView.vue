<script setup>
import { onMounted, ref, computed } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const assets = ref([]);
const showModal = ref(false);
const isEditMode = ref(false);
const editingAssetId = ref(null);
const error = ref("");
const searchQuery = ref("");
const filterStatus = ref("all");
const isLoading = ref(false);

const sites = ref([]);

const form = ref({
  name: "",
  type: "",
  serial: "",
  ip_address: "",
  site_id: null,
  status: "",
  assigned_to: "",
});

const getToken = () => localStorage.getItem("token");
const role = computed(() => localStorage.getItem("role") || "");

const statusConfig = {
  Actif: { color: "#059669", bg: "#d1fae5", dot: "#10b981" },
  Inactif: { color: "#6b7280", bg: "#f3f4f6", dot: "#9ca3af" },
  "En maintenance": { color: "#d97706", bg: "#fef3c7", dot: "#f59e0b" },
  Défaillant: { color: "#dc2626", bg: "#fee2e2", dot: "#ef4444" },
};

const statusOptions = ["Actif", "Inactif", "En maintenance", "Défaillant"];
const typeOptions = ["Serveur", "Poste de travail", "Switch", "Routeur", "Imprimante", "NAS", "Firewall", "Autre"];

const getStatusStyle = (status) => {
  const cfg = statusConfig[status];
  if (!cfg) return {};
  return {
    color: cfg.color,
    backgroundColor: cfg.bg,
    "--dot": cfg.dot,
  };
};

const filteredAssets = computed(() => {
  let list = assets.value;
  if (filterStatus.value !== "all") list = list.filter((a) => a.status === filterStatus.value);
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    list = list.filter(
      (a) =>
        a.name?.toLowerCase().includes(q) ||
        a.type?.toLowerCase().includes(q) ||
        a.ip_address?.toLowerCase().includes(q) ||
        a.site?.toLowerCase().includes(q) ||
        a.assigned_to?.toLowerCase().includes(q)
    );
  }
  return list;
});

const stats = computed(() => ({
  total: assets.value.length,
  actif: assets.value.filter((a) => a.status === "Actif").length,
  maintenance: assets.value.filter((a) => a.status === "En maintenance").length,
  defaillant: assets.value.filter((a) => a.status === "Défaillant").length,
}));

const resetForm = () => {
  form.value = { name: "", type: "", serial: "", ip_address: "", site_id: null, status: "", assigned_to: "" };
  isEditMode.value = false;
  editingAssetId.value = null;
  error.value = "";
};

const fetchAssets = async () => {
  isLoading.value = true;
  try {
    const response = await api.get("/assets", { headers: { Authorization: `Bearer ${getToken()}` } });
    assets.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger les assets";
  } finally {
    isLoading.value = false;
  }
};

const openCreateModal = () => { resetForm(); showModal.value = true; };

const openEditModal = (asset) => {
  isEditMode.value = true;
  editingAssetId.value = asset.ID;
  form.value = { name: asset.name, type: asset.type, serial: asset.serial, ip_address: asset.ip_address, site_id: asset.site_id, status: asset.status, assigned_to: asset.assigned_to };
  showModal.value = true;
};

const saveAsset = async () => {
  try {
    if (isEditMode.value) {
      await api.put(`/assets/${editingAssetId.value}`, form.value, { headers: { Authorization: `Bearer ${getToken()}` } });
    } else {
      await api.post("/assets", form.value, { headers: { Authorization: `Bearer ${getToken()}` } });
    }
    showModal.value = false;
    resetForm();
    fetchAssets();
  } catch (err) {
    console.error(err);
    error.value = "Erreur lors de l'enregistrement";
  }
};

const deleteAsset = async (asset) => {
  if (!confirm(`Supprimer l'asset « ${asset.name} » ?`)) return;
  try {
    await api.delete(`/assets/${asset.ID}`, { headers: { Authorization: `Bearer ${getToken()}` } });
    fetchAssets();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de supprimer cet asset";
  }
};

const typeIcon = (type) => {
  const icons = { Serveur: "🖥", "Poste de travail": "💻", Switch: "🔀", Routeur: "📡", Imprimante: "🖨", NAS: "💾", Firewall: "🛡", Autre: "📦" };
  return icons[type] || "📦";
};

const fetchSites = async () => {
  try {
    const token = getToken();

    const response = await api.get("/sites", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    sites.value = response.data;
  } catch (err) {
    console.error(err);
  }
};

onMounted(() => {
  fetchAssets();
  fetchSites();
});
</script>

<template>
  <DashboardLayout>
    <div class="assets-page">

      <!-- Header -->
      <div class="page-header">
        <div>
          <p class="page-eyebrow">Inventaire</p>
          <h1 class="page-title">Assets</h1>
        </div>
        <button class="btn-primary" @click="openCreateModal">
          <span class="btn-icon">＋</span>
          Ajouter un asset
        </button>
      </div>

      <!-- Stats bar -->
      <div class="stats-bar">
        <div class="stat-card">
          <span class="stat-value">{{ stats.total }}</span>
          <span class="stat-label">Total</span>
        </div>
        <div class="stat-card stat-green">
          <span class="stat-value">{{ stats.actif }}</span>
          <span class="stat-label">Actifs</span>
        </div>
        <div class="stat-card stat-amber">
          <span class="stat-value">{{ stats.maintenance }}</span>
          <span class="stat-label">Maintenance</span>
        </div>
        <div class="stat-card stat-red">
          <span class="stat-value">{{ stats.defaillant }}</span>
          <span class="stat-label">Défaillants</span>
        </div>
      </div>

      <!-- Toolbar -->
      <div class="toolbar">
        <div class="search-wrap">
          <span class="search-icon">🔍</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Rechercher par nom, type, IP, site…"
            class="search-input"
          />
        </div>
        <div class="filter-tabs">
          <button
            v-for="opt in ['all', ...statusOptions]"
            :key="opt"
            class="filter-tab"
            :class="{ active: filterStatus === opt }"
            @click="filterStatus = opt"
          >
            {{ opt === 'all' ? 'Tous' : opt }}
          </button>
        </div>
      </div>

      <!-- Error -->
      <div v-if="error" class="alert-error">
        <span>⚠</span> {{ error }}
      </div>

      <!-- Table -->
      <div class="table-card">
        <div v-if="isLoading" class="table-loading">
          <span class="spinner"></span> Chargement…
        </div>

        <table v-else class="assets-table">
          <thead>
            <tr>
              <th>Asset</th>
              <th>Type</th>
              <th>Adresse IP</th>
              <th>N° Série</th>
              <th>Site</th>
              <th>Assigné à</th>
              <th>Statut</th>
              <th class="col-actions"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="asset in filteredAssets" :key="asset.ID" class="table-row">
              <td class="cell-name">
                <span class="asset-icon">{{ typeIcon(asset.type) }}</span>
                <span class="asset-name">{{ asset.name }}</span>
              </td>
              <td>
                <span class="type-chip">{{ asset.type }}</span>
              </td>
              <td>
                <code class="ip-code">{{ asset.ip_address || '—' }}</code>
              </td>
              <td class="cell-muted">{{ asset.serial || '—' }}</td>
              <td>
                <span v-if="asset.site" class="site-tag">
                  <span class="site-dot">●</span>{{ asset.site?.name || "—" }}
                </span>
                <span v-else class="cell-muted">—</span>
              </td>
              <td>
                <div v-if="asset.assigned_to" class="assigned-user">
                  <span class="user-avatar">{{ asset.assigned_to.charAt(0).toUpperCase() }}</span>
                  <span>{{ asset.assigned_to }}</span>
                </div>
                <span v-else class="cell-muted">Non assigné</span>
              </td>
              <td>
                <span class="status-badge" :style="getStatusStyle(asset.status)">
                  <span class="status-dot"></span>
                  {{ asset.status }}
                </span>
              </td>
              <td class="cell-actions">
                <button class="action-btn edit-btn" @click="openEditModal(asset)" title="Modifier">
                  ✏
                </button>
                <button
                  v-if="role === 'Admin'"
                  class="action-btn delete-btn"
                  @click="deleteAsset(asset)"
                  title="Supprimer"
                >
                  🗑
                </button>
              </td>
            </tr>

            <tr v-if="filteredAssets.length === 0 && !isLoading">
              <td colspan="8" class="empty-state">
                <span class="empty-icon">📭</span>
                <p>Aucun asset trouvé</p>
                <small>Essayez de modifier vos filtres ou ajoutez un nouvel asset.</small>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="filteredAssets.length > 0" class="table-footer">
          {{ filteredAssets.length }} asset{{ filteredAssets.length > 1 ? 's' : '' }} affiché{{ filteredAssets.length > 1 ? 's' : '' }}
          <span v-if="assets.length !== filteredAssets.length"> sur {{ assets.length }} au total</span>
        </div>
      </div>
    </div>
  </DashboardLayout>

  <!-- Modal -->
  <Teleport to="body">
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <div class="modal-header">
          <div>
            <p class="modal-eyebrow">{{ isEditMode ? 'Modification' : 'Création' }}</p>
            <h2 class="modal-title">{{ isEditMode ? 'Modifier un asset' : 'Ajouter un asset' }}</h2>
          </div>
          <button class="modal-close" @click="showModal = false">✕</button>
        </div>

        <div v-if="error" class="alert-error modal-alert">
          <span>⚠</span> {{ error }}
        </div>

        <div class="modal-body">
          <div class="form-section">
            <p class="form-section-label">Identification</p>
            <div class="form-row">
              <div class="form-group">
                <label>Nom de l'asset <span class="required">*</span></label>
                <input v-model="form.name" type="text" placeholder="ex: SRV-PROD-01" />
              </div>
              <div class="form-group">
                <label>Type <span class="required">*</span></label>
                <select v-model="form.type">
                  <option value="" disabled>Sélectionner…</option>
                  <option v-for="t in typeOptions" :key="t" :value="t">{{ t }}</option>
                </select>
              </div>
            </div>
            <div class="form-group">
              <label>Numéro de série</label>
              <input v-model="form.serial" type="text" placeholder="ex: SN-XXXX-XXXX" />
            </div>
          </div>

          <div class="form-section">
            <p class="form-section-label">Réseau & localisation</p>
            <div class="form-row">
              <div class="form-group">
                <label>Adresse IP</label>
                <input v-model="form.ip_address" type="text" placeholder="192.168.1.1" />
              </div>
              <div class="form-group">
                <label>Site</label>
                <select v-model="form.site_id">
                  <option :value="null">
                    Sélectionner un site
                  </option>

                  <option
                      v-for="site in sites"
                      :key="site.ID"
                      :value="site.ID"
                  >
                    {{ site.name }}
                  </option>
                </select>
              </div>
            </div>
          </div>

          <div class="form-section">
            <p class="form-section-label">Affectation</p>
            <div class="form-row">
              <div class="form-group">
                <label>Statut</label>
                <select v-model="form.status">
                  <option value="" disabled>Sélectionner…</option>
                  <option v-for="s in statusOptions" :key="s" :value="s">{{ s }}</option>
                </select>
              </div>
              <div class="form-group">
                <label>Assigné à</label>
                <input v-model="form.assigned_to" type="text" placeholder="Nom de l'utilisateur" />
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn-ghost" @click="showModal = false">Annuler</button>
          <button class="btn-primary" @click="saveAsset">
            {{ isEditMode ? 'Enregistrer les modifications' : 'Créer l\'asset' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
/* ─── Page ─────────────────────────────────────────────── */
.assets-page {
  padding: 2rem;
  max-width: 1200px;
  font-family: 'Inter', 'Segoe UI', system-ui, sans-serif;
}

/* ─── Header ────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  margin-bottom: 1.75rem;
}
.page-eyebrow {
  font-size: 0.75rem;
  font-weight: 500;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #6366f1;
  margin: 0 0 0.25rem;
}
.page-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
  letter-spacing: -0.02em;
}

/* ─── Stats ─────────────────────────────────────────────── */
.stats-bar {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
  margin-bottom: 1.5rem;
}
.stat-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 1rem 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}
.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: #0f172a;
  line-height: 1;
}
.stat-label {
  font-size: 0.8rem;
  color: #64748b;
  font-weight: 500;
}
.stat-card.stat-green .stat-value { color: #059669; }
.stat-card.stat-amber .stat-value { color: #d97706; }
.stat-card.stat-red .stat-value { color: #dc2626; }
.stat-card.stat-green { border-left: 3px solid #10b981; }
.stat-card.stat-amber { border-left: 3px solid #f59e0b; }
.stat-card.stat-red { border-left: 3px solid #ef4444; }

/* ─── Toolbar ───────────────────────────────────────────── */
.toolbar {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.25rem;
  flex-wrap: wrap;
}
.search-wrap {
  flex: 1;
  min-width: 240px;
  position: relative;
}
.search-icon {
  position: absolute;
  left: 0.875rem;
  top: 50%;
  transform: translateY(-50%);
  font-size: 0.9rem;
  pointer-events: none;
}
.search-input {
  width: 100%;
  padding: 0.625rem 0.875rem 0.625rem 2.5rem;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  font-size: 0.875rem;
  background: #fff;
  color: #0f172a;
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
  box-sizing: border-box;
}
.search-input:focus {
  border-color: #818cf8;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.12);
}
.filter-tabs {
  display: flex;
  gap: 0.375rem;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 0.25rem;
}
.filter-tab {
  border: none;
  background: transparent;
  padding: 0.375rem 0.875rem;
  border-radius: 7px;
  font-size: 0.8rem;
  font-weight: 500;
  color: #64748b;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
  white-space: nowrap;
}
.filter-tab.active,
.filter-tab:hover {
  background: #fff;
  color: #0f172a;
  box-shadow: 0 1px 3px rgba(0,0,0,0.08);
}
.filter-tab.active {
  color: #6366f1;
  font-weight: 600;
}

/* ─── Alert ─────────────────────────────────────────────── */
.alert-error {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #b91c1c;
  padding: 0.75rem 1rem;
  border-radius: 10px;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}
.modal-alert { margin-bottom: 0; border-radius: 0; border-left: none; border-right: none; border-top: none; }

/* ─── Table card ─────────────────────────────────────────── */
.table-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  overflow: hidden;
}
.table-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  padding: 3rem;
  color: #64748b;
  font-size: 0.9rem;
}
.spinner {
  width: 18px;
  height: 18px;
  border: 2px solid #e2e8f0;
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  display: inline-block;
}
@keyframes spin { to { transform: rotate(360deg); } }

.assets-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}
.assets-table thead tr {
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}
.assets-table th {
  padding: 0.75rem 1rem;
  text-align: left;
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  white-space: nowrap;
}
.col-actions { width: 80px; }

.table-row {
  border-bottom: 1px solid #f1f5f9;
  transition: background 0.1s;
}
.table-row:last-child { border-bottom: none; }
.table-row:hover { background: #fafbff; }

.assets-table td {
  padding: 0.875rem 1rem;
  color: #1e293b;
  vertical-align: middle;
}

/* ─── Cell variants ─────────────────────────────────────── */
.cell-name {
  display: flex;
  align-items: center;
  gap: 0.625rem;
}
.asset-icon {
  font-size: 1.125rem;
  width: 32px;
  height: 32px;
  background: #f1f5f9;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.asset-name {
  font-weight: 600;
  color: #0f172a;
}
.cell-muted { color: #94a3b8; }

.type-chip {
  background: #f1f5f9;
  color: #475569;
  font-size: 0.75rem;
  font-weight: 500;
  padding: 0.25rem 0.625rem;
  border-radius: 6px;
}
.ip-code {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 0.8rem;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  padding: 0.2rem 0.5rem;
  border-radius: 5px;
  color: #334155;
}
.site-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.8rem;
  color: #334155;
}
.site-dot {
  font-size: 0.5rem;
  color: #94a3b8;
}
.assigned-user {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
}
.user-avatar {
  width: 26px;
  height: 26px;
  background: #e0e7ff;
  color: #6366f1;
  border-radius: 50%;
  font-size: 0.7rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

/* ─── Status badge ──────────────────────────────────────── */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.3rem 0.7rem;
  border-radius: 999px;
  white-space: nowrap;
}
.status-dot {
  width: 6px;
  height: 6px;
  background: var(--dot, currentColor);
  border-radius: 50%;
  flex-shrink: 0;
}

/* ─── Actions ───────────────────────────────────────────── */
.cell-actions {
  display: flex;
  gap: 0.25rem;
  justify-content: flex-end;
}
.action-btn {
  border: none;
  background: transparent;
  padding: 0.375rem;
  border-radius: 7px;
  cursor: pointer;
  font-size: 0.9rem;
  opacity: 0.5;
  transition: opacity 0.15s, background 0.15s;
}
.action-btn:hover { opacity: 1; }
.edit-btn:hover { background: #eff6ff; }
.delete-btn:hover { background: #fef2f2; }

/* ─── Empty state ───────────────────────────────────────── */
.empty-state {
  text-align: center;
  padding: 3.5rem 1rem !important;
  color: #64748b;
}
.empty-icon { font-size: 2.5rem; display: block; margin-bottom: 0.75rem; }
.empty-state p { font-size: 1rem; font-weight: 500; color: #334155; margin: 0 0 0.25rem; }
.empty-state small { font-size: 0.8rem; color: #94a3b8; }

/* ─── Table footer ──────────────────────────────────────── */
.table-footer {
  padding: 0.75rem 1rem;
  border-top: 1px solid #f1f5f9;
  font-size: 0.8rem;
  color: #94a3b8;
  background: #fafbff;
}

/* ─── Buttons ───────────────────────────────────────────── */
.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: #6366f1;
  color: #fff;
  border: none;
  padding: 0.625rem 1.125rem;
  border-radius: 10px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s, box-shadow 0.15s, transform 0.1s;
}
.btn-primary:hover { background: #4f46e5; box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3); }
.btn-primary:active { transform: scale(0.98); }
.btn-icon { font-size: 1rem; line-height: 1; }
.btn-ghost {
  background: transparent;
  border: 1px solid #e2e8f0;
  color: #374151;
  padding: 0.625rem 1.125rem;
  border-radius: 10px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.15s;
}
.btn-ghost:hover { background: #f8fafc; }

/* ─── Modal ─────────────────────────────────────────────── */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  backdrop-filter: blur(2px);
}
.modal {
  background: #fff;
  border-radius: 18px;
  width: 100%;
  max-width: 580px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 24px 64px rgba(0,0,0,0.15);
}
.modal-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 1.5rem 1.75rem 1.25rem;
  border-bottom: 1px solid #f1f5f9;
}
.modal-eyebrow {
  font-size: 0.7rem;
  font-weight: 600;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: #6366f1;
  margin: 0 0 0.2rem;
}
.modal-title {
  font-size: 1.2rem;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
}
.modal-close {
  border: none;
  background: #f1f5f9;
  color: #64748b;
  width: 30px;
  height: 30px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s;
  flex-shrink: 0;
}
.modal-close:hover { background: #e2e8f0; color: #0f172a; }

.modal-body {
  padding: 1.25rem 1.75rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}
.form-section-label {
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #94a3b8;
  margin: 0 0 0.875rem;
}
.form-section {
  display: flex;
  flex-direction: column;
}
.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}
.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
  margin-bottom: 0.875rem;
}
.form-group:last-child { margin-bottom: 0; }
.form-group label {
  font-size: 0.8rem;
  font-weight: 600;
  color: #374151;
}
.required { color: #ef4444; }
.form-group input,
.form-group select {
  border: 1px solid #e2e8f0;
  border-radius: 9px;
  padding: 0.6rem 0.875rem;
  font-size: 0.875rem;
  color: #0f172a;
  background: #fff;
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
}
.form-group input:focus,
.form-group select:focus {
  border-color: #818cf8;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.12);
}
.form-group input::placeholder { color: #c0ccda; }

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1.25rem 1.75rem;
  border-top: 1px solid #f1f5f9;
  background: #fafbff;
  border-radius: 0 0 18px 18px;
}
</style>