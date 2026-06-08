<script setup>
import { onMounted, ref, computed } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

import DataTable from "datatables.net-vue3";
import DataTablesCore from "datatables.net-dt";

import "datatables.net-dt/css/dataTables.dataTables.css";

DataTable.use(DataTablesCore);

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

const impact = ref(null);
const showImpactModal = ref(false);


const showImpact = async (asset) => {
  try {
    const response = await api.get(
        `/cmdb/impact/${asset.id || asset.ID}`,
        {
          headers: {
            Authorization: `Bearer ${getToken()}`
          }
        }
    );

    impact.value = response.data;
    showImpactModal.value = true;
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger l'analyse d'impact";
  }
};

const getToken = () => localStorage.getItem("token");
const role = computed(() => localStorage.getItem("role") || "");

const statusConfig = {
  Active: { color: "#34d399", bg: "rgba(52, 211, 153, 0.1)", dot: "#10b981" },
  Maintenance: { color: "#fbbf24", bg: "rgba(251, 191, 36, 0.1)", dot: "#f59e0b" },
  Retired: { color: "#9ca3af", bg: "rgba(156, 163, 175, 0.1)", dot: "#9ca3af" },
};

const statusOptions = ["Active", "Maintenance", "Retired"];
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

const columns = [
  { data: "name", title: "Asset", defaultContent: "Sans nom" },
  { data: "type", title: "Type", defaultContent: "Autre" },
  { data: "ip_address", title: "Adresse IP", defaultContent: "—" },
  { data: "serial", title: "N° Série", defaultContent: "—" },
  { data: "site.name", title: "Site", defaultContent: "—" },
  { data: "assigned_to", title: "Assigné à", defaultContent: "" },
  { data: "status", title: "Statut", defaultContent: "Active" },
  { data: null, title: "", orderable: false, className: "col-actions" }
];

const dtOptions = {
  responsive: true,
  pageLength: 10,
  dom: 't<"dt-footer flex items-center justify-between p-4"ip>',
  language: {
    info: "Affichage de _START_ à _END_ sur _TOTAL_ assets",
    infoEmpty: "Aucun équipement enregistré",
    infoFiltered: "(filtré depuis _MAX_ éléments)",
    zeroRecords: "Aucun asset trouvé",
    paginate: {
      next: "➔",
      previous: " Zar "
    }
  }
};

const fetchAssets = async () => {
  isLoading.value = true;
  try {
    const response = await api.get("/assets", { headers: { Authorization: `Bearer ${getToken()}` } });
    assets.value = response.data || [];
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger les actifs du parc informatique";
  } finally {
    isLoading.value = false;
  }
};

const filteredAssets = computed(() => {
  let list = assets.value;
  if (filterStatus.value !== "all") {
    list = list.filter((a) => a.status === filterStatus.value);
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    list = list.filter(
        (a) =>
            a.name?.toLowerCase().includes(q) ||
            a.type?.toLowerCase().includes(q) ||
            a.ip_address?.toLowerCase().includes(q) ||
            a.site?.name?.toLowerCase().includes(q) ||
            a.assigned_to?.toLowerCase().includes(q)
    );
  }
  return list;
});

const stats = computed(() => ({
  total: assets.value.length,
  actif: assets.value.filter((a) => a.status === "Active").length,
  maintenance: assets.value.filter((a) => a.status === "Maintenance").length,
  defaillant: assets.value.filter((a) => a.status === "Retired").length,
}));

const resetForm = () => {
  form.value = { name: "", type: "", serial: "", ip_address: "", site_id: null, status: "", assigned_to: "" };
  isEditMode.value = false;
  editingAssetId.value = null;
  error.value = "";
};

const openCreateModal = () => { resetForm(); showModal.value = true; };

const openEditModal = (asset) => {
  isEditMode.value = true;
  editingAssetId.value = asset.ID;
  form.value = {
    name: asset.name || "",
    type: asset.type || "",
    serial: asset.serial || "",
    ip_address: asset.ip_address || "",
    site_id: asset.site_id || null,
    status: asset.status || "",
    assigned_to: asset.assigned_to || ""
  };
  showModal.value = true;
};

const saveAsset = async () => {
  try {
    error.value = "";
    const payload = {
      ...form.value,
      site_id: form.value.site_id ? Number(form.value.site_id) : null,
    };

    if (isEditMode.value) {
      await api.put(`/assets/${editingAssetId.value}`, payload, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    } else {
      await api.post("/assets", payload, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    }

    showModal.value = false;
    resetForm();
    await fetchAssets();
  } catch (err) {
    console.error(err);
    error.value = "Erreur lors de la synchronisation de l'asset";
  }
};

const deleteAsset = async (asset) => {
  if (!confirm(`Supprimer l'asset « ${asset.name} » ?`)) return;
  try {
    await api.delete(`/assets/${asset.ID}`, { headers: { Authorization: `Bearer ${getToken()}` } });
    await fetchAssets();
  } catch (err) {
    console.error(err);
    error.value = "Impossible d'isoler ou de supprimer cet asset";
  }
};

const typeIcon = (type) => {
  const icons = { Serveur: "🖥", "Poste de travail": "💻", Switch: "🔀", Routeur: "📡", Imprimante: "🖨", NAS: "💾", Firewall: "🛡", Autre: "📦" };
  return icons[type] || "📦";
};

const fetchSites = async () => {
  try {
    const response = await api.get("/sites", { headers: { Authorization: `Bearer ${getToken()}` } });
    sites.value = response.data || [];
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

      <div class="page-header">
        <div>
          <p class="page-eyebrow">Inventaire du Parc</p>
          <h1 class="page-title">Assets</h1>
        </div>
        <button class="btn-primary" @click="openCreateModal">
          <span class="btn-icon">＋</span>
          Ajouter un asset
        </button>
      </div>

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

      <div class="toolbar">
        <div class="search-wrap">
          <span class="search-icon">🔍</span>
          <input
              v-model="searchQuery"
              type="text"
              placeholder="Rechercher par nom, type, IP, site ou assignation…"
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

      <div v-if="error" class="alert-error">
        <span>⚠</span> {{ error }}
      </div>

      <div class="table-card assets-dt-wrapper">
        <div v-if="isLoading" class="table-loading">
          <span class="spinner"></span> Chargement des assets...
        </div>

        <DataTable
            v-else
            :data="filteredAssets"
            :columns="columns"
            class="assets-table w-full text-left border-collapse"
            :options="dtOptions"
        >
          <template #column-name="{ rowData }">
            <div class="cell-name">
              <span class="asset-icon">{{ typeIcon(rowData.type) }}</span>
              <span class="asset-name">{{ rowData.name }}</span>
            </div>
          </template>

          <template #column-type="{ cellData }">
            <span class="type-chip">{{ cellData }}</span>
          </template>

          <template #column-ip_address="{ cellData }">
            <code v-if="cellData && cellData !== '—'" class="ip-code">{{ cellData }}</code>
            <span v-else class="cell-muted">—</span>
          </template>

          <template #column-site-name="{ cellData }">
            <span v-if="cellData && cellData !== '—'" class="site-tag">
              <span class="site-dot">●</span>{{ cellData }}
            </span>
            <span v-else class="cell-muted">—</span>
          </template>

          <template #column-assigned_to="{ cellData }">
            <div v-if="cellData" class="assigned-user">
              <span class="user-avatar">{{ cellData.charAt(0).toUpperCase() }}</span>
              <span>{{ cellData }}</span>
            </div>
            <span v-else class="cell-muted">Non assigné</span>
          </template>

          <template #column-status="{ cellData }">
            <span class="status-badge" :style="getStatusStyle(cellData)">
              <span class="status-dot"></span>
              {{ cellData }}
            </span>
          </template>

          <template #column-7="{ rowData }">
            <div class="cell-actions">
              <button
                  class="action-btn impact-btn"
                  @click="showImpact(rowData)"
                  title="Analyse d'impact"
              >
                🧩
              </button>

              <button
                  class="action-btn edit-btn"
                  @click="openEditModal(rowData)"
                  title="Modifier"
              >
                ✏
              </button>

              <button
                  v-if="role === 'Admin'"
                  class="action-btn delete-btn"
                  @click="deleteAsset(rowData)"
                  title="Supprimer"
              >
                🗑
              </button>
            </div>
          </template>
        </DataTable>

        <div v-if="filteredAssets.length > 0 && !isLoading" class="table-footer">
          {{ filteredAssets.length }} asset{{ filteredAssets.length > 1 ? 's' : '' }} affiché{{ filteredAssets.length > 1 ? 's' : '' }}
          <span v-if="assets.length !== filteredAssets.length"> sur {{ assets.length }} au total</span>
        </div>
      </div>
    </div>
  </DashboardLayout>

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
                  <option :value="null">Sélectionner un site</option>
                  <option v-for="site in sites" :key="site.ID" :value="site.ID">
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

  <Teleport to="body">
    <div
        v-if="showImpactModal"
        class="fixed inset-0 bg-black/70 flex items-center justify-center z-50"
        @click.self="showImpactModal = false"
    >
      <div class="bg-slate-900 border border-slate-700 w-full max-w-2xl rounded-2xl p-6 text-white">
        <div class="flex justify-between items-start mb-6">
          <div>
            <p class="text-xs uppercase tracking-wider text-indigo-400 font-semibold">
              CMDB Impact Analysis
            </p>
            <h2 class="text-xl font-bold">
              {{ impact?.asset?.name || "Asset" }}
            </h2>
          </div>

          <button
              @click="showImpactModal = false"
              class="text-slate-400 hover:text-white"
          >
            ✕
          </button>
        </div>
        <div class="bg-slate-800 rounded-xl p-4 border border-slate-700 mb-4">
          <h3 class="text-sm font-semibold text-red-400 mb-3">
            Business Services impactés
          </h3>

          <ul v-if="impact?.business_services?.length">
            <li
                v-for="service in impact.business_services"
                :key="service.ID"
                class="py-2 border-b border-slate-700 last:border-0"
            >
              <div class="font-medium">
                {{ service.name }}
              </div>

              <div class="text-xs text-slate-400">
                Criticité : {{ service.criticality }}
              </div>
            </li>
          </ul>

          <p
              v-else
              class="text-sm text-slate-400"
          >
            Aucun service métier impacté.
          </p>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="bg-slate-800 rounded-xl p-4 border border-slate-700">
            <h3 class="text-sm font-semibold text-indigo-400 mb-3">
              Applications impactées
            </h3>

            <ul v-if="impact?.applications?.length">
              <li
                  v-for="app in impact.applications"
                  :key="app.ID"
                  class="py-2 border-b border-slate-700 last:border-0"
              >
                <div class="font-medium">{{ app.name }}</div>
                <div class="text-xs text-slate-400">
                  {{ app.criticality }} · {{ app.status }}
                </div>
              </li>
            </ul>

            <p v-else class="text-sm text-slate-400">
              Aucune application impactée.
            </p>
          </div>

          <div class="bg-slate-800 rounded-xl p-4 border border-slate-700">
            <h3 class="text-sm font-semibold text-emerald-400 mb-3">
              Bases de données impactées
            </h3>

            <ul v-if="impact?.databases?.length">
              <li
                  v-for="db in impact.databases"
                  :key="db.ID"
                  class="py-2 border-b border-slate-700 last:border-0"
              >
                <div class="font-medium">{{ db.name }}</div>
                <div class="text-xs text-slate-400">
                  {{ db.engine }} {{ db.version }} · {{ db.environment }}
                </div>
              </li>
            </ul>

            <p v-else class="text-sm text-slate-400">
              Aucune base impactée.
            </p>
          </div>
        </div>

        <div class="mt-6 flex justify-end">
          <button
              @click="showImpactModal = false"
              class="px-4 py-2 bg-slate-700 hover:bg-slate-600 rounded-xl text-sm"
          >
            Fermer
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>



<style scoped>
/* Refonte complète pour le Dark Mode (Palette Slate-900 / Slate-800) */
.assets-page {
  padding: 2rem;
  max-width: 1200px;
  font-family: 'Inter', 'Segoe UI', system-ui, sans-serif;
  background-color: #0f172a; /* Fond noir/sombre profond */
  min-height: 100vh;
  color: #f1f5f9;
}
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
  color: #818cf8;
  margin: 0 0 0.25rem;
}
.page-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #f8fafc;
  margin: 0;
  letter-spacing: -0.02em;
}
.stats-bar {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
  margin-bottom: 1.5rem;
}
.stat-card {
  background: #1e293b; /* Cartes sombres */
  border: 1px solid #334155;
  border-radius: 12px;
  padding: 1rem 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}
.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: #f8fafc;
  line-height: 1;
}
.stat-label {
  font-size: 0.8rem;
  color: #94a3b8;
  font-weight: 500;
}
.stat-card.stat-green .stat-value { color: #34d399; }
.stat-card.stat-amber .stat-value { color: #fbbf24; }
.stat-card.stat-red .stat-value { color: #f87171; }
.stat-card.stat-green { border-left: 3px solid #10b981; }
.stat-card.stat-amber { border-left: 3px solid #f59e0b; }
.stat-card.stat-red { border-left: 3px solid #ef4444; }

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
  color: #64748b;
}
.search-input {
  width: 100%;
  padding: 0.625rem 0.875rem 0.625rem 2.5rem;
  border: 1px solid #334155;
  border-radius: 10px;
  font-size: 0.875rem;
  background: #1e293b;
  color: #f8fafc;
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
  box-sizing: border-box;
}
.search-input:focus {
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.25);
}
.filter-tabs {
  display: flex;
  gap: 0.375rem;
  background: #1e293b;
  border: 1px solid #334155;
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
  color: #94a3b8;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
  white-space: nowrap;
}
.filter-tab:hover {
  background: #273549;
  color: #f8fafc;
}
.filter-tab.active {
  background: #0f172a;
  color: #818cf8;
  box-shadow: 0 1px 3px rgba(0,0,0,0.3);
  font-weight: 600;
}

.alert-error {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid #ef4444;
  color: #fca5a5;
  padding: 0.75rem 1rem;
  border-radius: 10px;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}
.modal-alert { margin-bottom: 0; border-radius: 0; border-left: none; border-right: none; border-top: none; }

.table-card {
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 16px;
  overflow: hidden;
}
.table-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  padding: 3rem;
  color: #94a3b8;
  font-size: 0.9rem;
}
.spinner {
  width: 18px;
  height: 18px;
  border: 2px solid #334155;
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  display: inline-block;
}
@keyframes spin { to { transform: rotate(360deg); } }

.cell-name {
  display: flex;
  align-items: center;
  gap: 0.625rem;
}
.asset-icon {
  font-size: 1.125rem;
  width: 32px;
  height: 32px;
  background: #2d3748;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.asset-name {
  font-weight: 600;
  color: #f8fafc;
}
.cell-muted { color: #64748b; }

.type-chip {
  background: #2d3748;
  color: #cbd5e1;
  font-size: 0.75rem;
  font-weight: 500;
  padding: 0.25rem 0.625rem;
  border-radius: 6px;
  display: inline-block;
}
.ip-code {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 0.8rem;
  background: #0f172a;
  border: 1px solid #334155;
  padding: 0.2rem 0.5rem;
  border-radius: 5px;
  color: #e2e8f0;
}
.site-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.8rem;
  color: #e2e8f0;
}
.site-dot {
  font-size: 0.5rem;
  color: #64748b;
}
.assigned-user {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  color: #e2e8f0;
}
.user-avatar {
  width: 26px;
  height: 26px;
  background: rgba(99, 102, 241, 0.2);
  color: #818cf8;
  border-radius: 50%;
  font-size: 0.7rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

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
  opacity: 0.6;
  transition: opacity 0.15s, background 0.15s;
}
.action-btn:hover { opacity: 1; }
.edit-btn:hover { background: #2d3748; }
.delete-btn:hover { background: rgba(239, 68, 68, 0.2); }

.table-footer {
  padding: 0.75rem 1rem;
  border-top: 1px solid #334155;
  font-size: 0.8rem;
  color: #64748b;
  background: #1e293b;
}

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
.btn-primary:hover { background: #4f46e5; box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4); }
.btn-primary:active { transform: scale(0.98); }
.btn-icon { font-size: 1rem; line-height: 1; }

.btn-ghost {
  background: transparent;
  border: 1px solid #334155;
  color: #cbd5e1;
  padding: 0.625rem 1.125rem;
  border-radius: 10px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.15s;
}
.btn-ghost:hover { background: #2d3748; }

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  backdrop-filter: blur(4px);
}
.modal {
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 18px;
  width: 100%;
  max-width: 580px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 24px 64px rgba(0,0,0,0.5);
}
.modal-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 1.5rem 1.75rem 1.25rem;
  border-bottom: 1px solid #334155;
}
.modal-eyebrow {
  font-size: 0.7rem;
  font-weight: 600;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: #818cf8;
  margin: 0 0 0.2rem;
}
.modal-title {
  font-size: 1.2rem;
  font-weight: 700;
  color: #f8fafc;
  margin: 0;
}
.modal-close {
  border: none;
  background: #2d3748;
  color: #94a3b8;
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
.modal-close:hover { background: #334155; color: #f8fafc; }

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
  color: #64748b;
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
  color: #cbd5e1;
}
.required { color: #f87171; }
.form-group input,
.form-group select {
  border: 1px solid #334155;
  border-radius: 9px;
  padding: 0.6rem 0.875rem;
  font-size: 0.875rem;
  color: #f8fafc;
  background: #0f172a;
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
}
.form-group input:focus,
.form-group select:focus {
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.25);
}
.form-group input::placeholder { color: #475569; }

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1.25rem 1.75rem;
  border-top: 1px solid #334155;
  background: #1e293b;
  border-radius: 0 0 18px 18px;
}
</style>

<style>
.impact-btn:hover {
  background: rgba(99, 102, 241, 0.15);
  color: #818cf8;
}
.assets-dt-wrapper .dataTables_wrapper {
  background: transparent !important;
  border: none !important;
  padding: 0 !important;
}

.assets-dt-wrapper table.dataTable {
  border-collapse: collapse !important;
  margin: 0 !important;
  background: transparent !important;
  width: 100% !important;
  border-bottom: none !important;
}

.assets-dt-wrapper table.dataTable thead tr {
  background: #1e293b !important;
  border-bottom: 1px solid #334155 !important;
}

.assets-dt-wrapper table.dataTable th {
  padding: 0.75rem 1rem !important;
  font-size: 0.75rem !important;
  font-weight: 600 !important;
  color: #94a3b8 !important;
  text-transform: uppercase !important;
  letter-spacing: 0.06em !important;
  border-bottom: 1px solid #334155 !important;
  position: relative;
  background-image: none !important;
}

.assets-dt-wrapper table.dataTable thead th.sorting::after,
.assets-dt-wrapper table.dataTable thead th.sorting_asc::after,
.assets-dt-wrapper table.dataTable thead th.sorting_desc::after {
  position: absolute;
  right: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  font-size: 0.65rem;
  opacity: 0.4;
  transition: opacity 0.15s, color 0.15s;
}

.assets-dt-wrapper table.dataTable thead th.sorting::after { content: "▲▼"; color: #64748b;}
.assets-dt-wrapper table.dataTable thead th.sorting_asc::after { content: "▲"; opacity: 1; color: #818cf8; }
.assets-dt-wrapper table.dataTable thead th.sorting_desc::after { content: "▼"; opacity: 1; color: #818cf8; }

.assets-dt-wrapper table.dataTable tbody td {
  padding: 0.875rem 1rem !important;
  vertical-align: middle !important;
  border-bottom: 1px solid #334155 !important;
  background: transparent !important;
  color: #e2e8f0;
}

.assets-dt-wrapper table.dataTable tbody tr {
  background-color: transparent !important;
}

.assets-dt-wrapper table.dataTable tbody tr:hover td {
  background-color: #243146 !important; /* Effet hover subtil en ligne */
}

.dt-footer {
  border-top: 1px solid #334155;
  background: #1e293b;
  color: #64748b;
  font-size: 0.8rem;
}

/* Pagination Dark Mode */
.dataTables_wrapper .dataTables_paginate {
  padding-top: 0 !important;
  margin: 0 !important;
}

.dataTables_wrapper .dataTables_paginate .paginate_button {
  background: #0f172a !important;
  border: 1px solid #334155 !important;
  color: #94a3b8 !important;
  border-radius: 8px !important;
  padding: 0.35rem 0.75rem !important;
  font-size: 0.8rem !important;
  font-weight: 500 !important;
  transition: all 0.15s ease !important;
  margin-left: 0.375rem !important;
}

.dataTables_wrapper .dataTables_paginate .paginate_button:hover {
  background: #2d3748 !important;
  border-color: #475569 !important;
  color: #f8fafc !important;
}

.dataTables_wrapper .dataTables_paginate .paginate_button.current,
.dataTables_wrapper .dataTables_paginate .paginate_button.current:hover {
  background: #6366f1 !important;
  color: white !important;
  border-color: #6366f1 !important;
  font-weight: 600 !important;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.4) !important;
}

.dataTables_wrapper .dataTables_paginate .paginate_button.disabled,
.dataTables_wrapper .dataTables_paginate .paginate_button.disabled:hover,
.dataTables_wrapper .dataTables_paginate .paginate_button.disabled:active {
  background: #1e293b !important;
  border-color: #334155 !important;
  color: #475569 !important;
  cursor: not-allowed !important;
  box-shadow: none !important;
}
</style>