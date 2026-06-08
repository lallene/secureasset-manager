<script setup>
import { onMounted, ref, computed } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

import DataTable from "datatables.net-vue3";
import DataTablesCore from "datatables.net-dt";

import "datatables.net-dt/css/dataTables.dataTables.css";

DataTable.use(DataTablesCore);

const databases = ref([]);
const sites = ref([]);

const showModal = ref(false);
const isEditMode = ref(false);
const editingDatabaseId = ref(null);

const error = ref("");
const isSubmitting = ref(false);
const searchQuery = ref("");
const filterEngine = ref("all");

const form = ref({
  name: "",
  engine: "PostgreSQL",
  version: "",
  environment: "Production",
  site_id: "",
});

const getToken = () => localStorage.getItem("token");
const role = computed(() => localStorage.getItem("role") || "");

const engineOptions = ["PostgreSQL", "MySQL", "SQL Server", "Oracle", "MongoDB", "MariaDB", "Redis", "Autre"];
const environmentOptions = ["Production", "Preproduction", "Staging", "Development", "Test"];

// Mapping sémantique aligné sur tes variables globales CSS
const engineConfig = {
  PostgreSQL:   { class: "p-ind",   label: "PostgreSQL" },
  MySQL:        { class: "p-amb",   label: "MySQL" },
  MariaDB:      { class: "p-amb",   label: "MariaDB" },
  "SQL Server": { class: "p-rose",  label: "SQL Server" },
  Oracle:       { class: "p-sl",    label: "Oracle" },
  MongoDB:      { class: "p-grn",   label: "MongoDB" },
  Redis:        { class: "p-rose",  label: "Redis" },
  Autre:        { class: "p-sl",    label: "Autre" }
};

const environmentConfig = {
  Production:    { class: "bdg-rose pulse-slow", label: "Production" },
  Preproduction: { class: "bdg-ind",             label: "Preprod" },
  Staging:       { class: "bdg-amb",             label: "Staging" },
  Development:   { class: "bdg-sl",              label: "Dev" },
  Test:          { class: "bdg-sl",              label: "Test" }
};

const columns = [
  { data: "name", title: "Instance / Base" },
  { data: "engine", title: "Moteur" },
  { data: "version", title: "Version" },
  { data: "environment", title: "Environnement" },
  { data: "site.name", title: "Site", defaultContent: "—" },
  { data: null, title: "Actions", orderable: false, className: "text-right" }
];

const dtOptions = {
  responsive: true,
  pageLength: 10,
  dom: 't<"dt-footer flex items-center justify-between p-4"ip>',
  language: {
    info: "Affichage de _START_ à _END_ sur _TOTAL_ instances",
    infoEmpty: "Aucune instance enregistrée",
    infoFiltered: "(filtré depuis _MAX_ éléments)",
    zeroRecords: "Aucune correspondance trouvée",
    paginate: {
      next: "➔",
      previous: " Zar "
    }
  }
};

const filteredDatabases = computed(() => {
  let list = databases.value;

  if (filterEngine.value !== "all") {
    list = list.filter((db) => db.engine === filterEngine.value);
  }

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    list = list.filter(
        (db) =>
            db.name?.toLowerCase().includes(q) ||
            db.engine?.toLowerCase().includes(q) ||
            db.version?.toLowerCase().includes(q) ||
            db.environment?.toLowerCase().includes(q) ||
            db.site?.name?.toLowerCase().includes(q)
    );
  }

  return list;
});

const stats = computed(() => ({
  total: databases.value.length,
  production: databases.value.filter((db) => db.environment === "Production").length,
  postgres: databases.value.filter((db) => db.engine === "PostgreSQL").length,
  mysql: databases.value.filter((db) => db.engine === "MySQL").length,
}));

const fetchDatabases = async () => {
  try {
    const response = await api.get("/cmdb/databases", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    databases.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger l'inventaire des bases de données";
  }
};

const fetchSites = async () => {
  try {
    const response = await api.get("/sites", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    sites.value = response.data;
  } catch (err) {}
};

const resetForm = () => {
  form.value = { name: "", engine: "PostgreSQL", version: "", environment: "Production", site_id: "" };
  isEditMode.value = false;
  editingDatabaseId.value = null;
  error.value = "";
};

const openCreateModal = () => { resetForm(); showModal.value = true; };

const openEditModal = (db) => {
  isEditMode.value = true;
  editingDatabaseId.value = db.ID;
  form.value = {
    name: db.name,
    engine: db.engine,
    version: db.version,
    environment: db.environment,
    site_id: db.site_id,
  };
  showModal.value = true;
};

const saveDatabase = async () => {
  try {
    isSubmitting.value = true;
    error.value = "";
    const payload = {
      ...form.value,
      site_id: form.value.site_id ? Number(form.value.site_id) : null,
    };

    if (isEditMode.value) {
      await api.put(`/cmdb/databases/${editingDatabaseId.value}`, payload, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    } else {
      await api.post("/cmdb/databases", payload, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    }
    showModal.value = false;
    resetForm();
    await fetchDatabases();
  } catch (err) {
    error.value = "Erreur lors de l'enregistrement de l'instance";
  } finally {
    isSubmitting.value = false;
  }
};

const deleteDatabase = async (db) => {
  if (!confirm(`Supprimer définitivement l'instance "${db.name}" ?`)) return;
  try {
    await api.delete(`/cmdb/databases/${db.ID}`, {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    await fetchDatabases();
  } catch (err) {
    error.value = "Rupture de contrainte : impossible de supprimer cet élément";
  }
};

onMounted(() => {
  fetchDatabases();
  fetchSites();
});
</script>

<template>
  <DashboardLayout>
    <div class="topbar flex items-center justify-between">
      <div>
        <span class="tb-env">Configuration Management Database (CMDB)</span>
        <h1>Bases de données</h1>
        <p class="text-sm text-slate-500 mt-1">
          Inventoriez les serveurs de bases de données, clusters, versions et environnements cibles.
        </p>
      </div>

      <button
          v-if="role === 'Admin'"
          @click="openCreateModal"
          class="btn-primary flex items-center gap-2"
      >
        <span>+</span> Ajouter une base
      </button>
    </div>

    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-muted flex items-center justify-center text-base w-11 h-11">🗄️</div>
        <div>
          <div class="kc-lbl text-xs">Total Instances</div>
          <div class="kc-val text-xl font-bold mt-0.5">{{ stats.total }}</div>
        </div>
      </div>
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-danger flex items-center justify-center text-base w-11 h-11">🚨</div>
        <div>
          <div class="kc-lbl text-xs">Bases de Production</div>
          <div class="kc-val text-xl font-bold mt-0.5" style="color: var(--c-rose-t)">{{ stats.production }}</div>
        </div>
      </div>
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-brand flex items-center justify-center text-base w-11 h-11">🐘</div>
        <div>
          <div class="kc-lbl text-xs">PostgreSQL</div>
          <div class="kc-val text-xl font-bold mt-0.5" style="color: var(--brand-light)">{{ stats.postgres }}</div>
        </div>
      </div>
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-warn flex items-center justify-center text-base w-11 h-11">🐬</div>
        <div>
          <div class="kc-lbl text-xs">MySQL / MariaDB</div>
          <div class="kc-val text-xl font-bold mt-0.5" style="color: var(--c-warn-t)">{{ stats.mysql }}</div>
        </div>
      </div>
    </div>

    <div class="fbar flex items-center justify-between gap-4 flex-wrap mb-6">
      <div class="search-wrap flex-1 relative">
        <input
            v-model="searchQuery"
            type="text"
            placeholder="Rechercher par nom d'instance, version de moteur, environnement, cluster..."
            class="search-input w-full pl-4 pr-10 py-2.5"
        />
        <span class="absolute right-4 top-3 text-sm">🔍</span>
      </div>

      <select v-model="filterEngine" class="field-input w-52">
        <option value="all">Tous les moteurs</option>
        <option v-for="engine in engineOptions" :key="engine" :value="engine">
          {{ engine }}
        </option>
      </select>
    </div>

    <div v-if="error" class="modal-error error-bar mb-6 flex items-center gap-2">
      <span>⚠️</span> {{ error }}
    </div>

    <div class="table-card overflow-hidden cmdb-dt-wrapper">
      <DataTable
          :data="filteredDatabases"
          :columns="columns"
          class="w-full text-left border-collapse"
          :options="dtOptions"
      >
        <template #column-name="{ cellData }">
          <div class="font-medium" style="color: var(--tx-primary)">{{ cellData }}</div>
        </template>

        <template #column-engine="{ cellData }">
          <span :class="['pill text-xs font-medium px-2.5 py-0.5 rounded-full border', engineConfig[cellData]?.class || 'p-sl']">
            {{ engineConfig[cellData]?.label || cellData }}
          </span>
        </template>

        <template #column-version="{ cellData }">
          <span class="px-1.5 py-0.5 rounded text-[11px] font-mono border" style="background: var(--bg-base); border-color: var(--bd-subtle); color: var(--tx-secondary)">
            {{ cellData || '—' }}
          </span>
        </template>

        <template #column-environment="{ cellData }">
          <span :class="['sc-bdg text-xs px-2.5 py-0.5 rounded-full border font-medium', environmentConfig[cellData]?.class || 'bdg-sl']">
            {{ environmentConfig[cellData]?.label || cellData }}
          </span>
        </template>

        <template #column-5="{ rowData }">
          <div class="flex justify-end gap-1.5">
            <button
                v-if="role === 'Admin'"
                class="action-btn action-edit text-xs"
                title="Modifier l'instance"
                @click="openEditModal(rowData)"
            >✏️</button>
            <button
                v-if="role === 'Admin'"
                class="action-btn action-delete text-xs"
                title="Purger du catalogue"
                @click="deleteDatabase(rowData)"
            >🗑️</button>
          </div>
        </template>
      </DataTable>
    </div>
  </DashboardLayout>

  <Teleport to="body">
    <div v-if="showModal" class="modal-overlay fixed inset-0 flex items-center justify-center p-4 z-50" @click.self="showModal = false">
      <div class="modal rounded-2xl w-full max-w-xl overflow-hidden shadow-2xl">

        <div class="modal-header p-5 flex justify-between items-center border-b" style="border-color: var(--bd-subtle)">
          <div>
            <span class="modal-eyebrow text-xs font-bold uppercase tracking-wider block mb-0.5" style="color: var(--brand-light)">Topologie Stockage</span>
            <h2 class="modal-title text-lg font-semibold m-0">
              {{ isEditMode ? "Modifier l'instance de données" : "Déclarer une ressource de stockage" }}
            </h2>
          </div>
          <button @click="showModal = false" class="modal-close p-2">✕</button>
        </div>

        <form @submit.prevent="saveDatabase">
          <div class="p-5 space-y-4 overflow-y-auto max-h-[70vh]">

            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium">Nom DNS ou Identifiant de l'instance <span style="color: var(--c-rose-t)">*</span></label>
              <input v-model="form.name" required type="text" class="field-input w-full" placeholder="Ex: prod-db-cluster-01" />
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-medium">Moteur SGBD <span style="color: var(--c-rose-t)">*</span></label>
                <select v-model="form.engine" required class="field-input w-full appearance-none cursor-pointer">
                  <option v-for="engine in engineOptions" :key="engine" :value="engine">{{ engine }}</option>
                </select>
              </div>

              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-medium">Version Tag</label>
                <input v-model="form.version" type="text" class="field-input w-full font-mono" placeholder="Ex: 16.2-alpine" />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-medium">Scope d'environnement <span style="color: var(--c-rose-t)">*</span></label>
                <select v-model="form.environment" required class="field-input w-full appearance-none cursor-pointer">
                  <option v-for="env in environmentOptions" :key="env" :value="env">{{ env }}</option>
                </select>
              </div>

              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-medium">Site d'attachement / Datacenter</label>
                <select v-model="form.site_id" class="field-input w-full appearance-none cursor-pointer">
                  <option value="">Sélectionner un cluster physique</option>
                  <option v-for="site in sites" :key="site.ID" :value="site.ID">{{ site.name }}</option>
                </select>
              </div>
            </div>

          </div>

          <div class="modal-footer p-4 flex justify-end gap-3 border-t" style="border-color: var(--bd-subtle); background: rgba(0,0,0,0.05)">
            <button type="button" class="btn-ghost" @click="showModal = false">Annuler</button>
            <button type="submit" :disabled="isSubmitting" class="btn-primary">
              {{ isEditMode ? "Enregistrer" : "Créer l'instance" }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.btn-primary, .btn-ghost {
  font-size: 14px;
  font-weight: 500;
  border-radius: var(--r-md);
  cursor: pointer;
  transition: all var(--t-fast);
}

.action-btn {
  padding: 6px 10px;
  font-weight: 500;
  border-radius: var(--r-sm);
  border: 1px solid transparent;
  cursor: pointer;
  transition: all var(--t-fast);
}

.error-bar {
  padding: 14px 16px;
  border-radius: var(--r-md);
  border: 1px solid;
  font-size: 14px;
}

.modal-header, .modal-footer {
  border-style: solid;
}

.pulse-slow {
  animation: pulseGlow 2s infinite ease-in-out;
}

@keyframes pulseGlow {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}
</style>

<style>
/* Encapsulation globale dédiée à DataTables pour match avec le framework d'UI sombre */
.cmdb-dt-wrapper table.dataTable {
  border-collapse: collapse !important;
  margin: 0 !important;
  background: transparent !important;
}

.cmdb-dt-wrapper table.dataTable thead th {
  background: rgba(0,0,0,0.15) !important;
  color: var(--tx-muted) !important;
  font-size: 11px !important;
  font-weight: 600 !important;
  text-transform: uppercase !important;
  letter-spacing: 0.05em !important;
  border-bottom: 1px solid var(--bd-subtle) !important;
  padding: 12px 16px !important;
}

.cmdb-dt-wrapper table.dataTable tbody td {
  padding: 14px 16px !important;
  border-bottom: 1px solid var(--bd-subtle) !important;
  background: transparent !important;
}

.cmdb-dt-wrapper table.dataTable tbody tr {
  background-color: transparent !important;
}

.cmdb-dt-wrapper table.dataTable tbody tr:hover td {
  background-color: rgba(255, 255, 255, 0.01) !important;
}

.dt-footer {
  border-top: 1px solid var(--bd-subtle);
  background: rgba(0, 0, 0, 0.08);
  color: var(--tx-muted);
  font-size: 12px;
}

.dataTables_wrapper .dataTables_paginate .paginate_button {
  background: var(--bg-base) !important;
  border: 1px solid var(--bd-subtle) !important;
  color: var(--tx-secondary) !important;
  border-radius: var(--r-sm) !important;
  padding: 4px 10px !important;
  font-size: 11px !important;
}

.dataTables_wrapper .dataTables_paginate .paginate_button.current {
  background: var(--brand) !important;
  color: white !important;
  border-color: var(--brand) !important;
}
</style>