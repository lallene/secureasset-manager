<script setup>
import { onMounted, ref, computed } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

import DataTable from "datatables.net-vue3";
import DataTablesCore from "datatables.net-dt";

import "datatables.net-dt/css/dataTables.dataTables.css";

DataTable.use(DataTablesCore);

const relations = ref([]);
const applications = ref([]);
const databases = ref([]);
const assets = ref([]);

const showModal = ref(false);
const isEditMode = ref(false);
const editingRelationId = ref(null);

const error = ref("");
const isSubmitting = ref(false);
const searchQuery = ref("");

const form = ref({
  source_type: "Application",
  source_id: "",
  relation_type: "hosted_on",
  target_type: "Asset",
  target_id: "",
});

const getToken = () => localStorage.getItem("token");
const role = computed(() => localStorage.getItem("role") || "");

const typeOptions = ["Application", "Asset", "Database"];
const relationOptions = ["hosted_on", "uses_database", "depends_on", "connects_to"];

const headers = () => ({
  headers: { Authorization: `Bearer ${getToken()}` },
});

// Mapping des libellés sémantiques pour les relations
const relationConfig = {
  hosted_on:     { class: "p-ind",  label: "hosted_on" },
  uses_database: { class: "p-rose", label: "uses_database" },
  depends_on:    { class: "p-amb",  label: "depends_on" },
  connects_to:   { class: "p-grn",  label: "connects_to" }
};

const columns = [
  { data: "source_id", title: "Élément Source" },
  { data: "relation_type", title: "Type de Liaison" },
  { data: "target_id", title: "Élément Cible" },
  { data: null, title: "Actions", orderable: false, className: "text-right" }
];

const dtOptions = {
  responsive: true,
  pageLength: 10,
  dom: 't<"dt-footer flex items-center justify-between p-4"ip>',
  language: {
    info: "Affichage de _START_ à _END_ sur _TOTAL_ liaisons topologiques",
    infoEmpty: "Aucune relation répertoriée",
    infoFiltered: "(filtré depuis _MAX_ éléments)",
    zeroRecords: "Aucune dépendance trouvée",
    paginate: {
      next: "➔",
      previous: " Zar "
    }
  }
};

const fetchRelations = async () => {
  try {
    const response = await api.get("/cmdb/relations", headers());
    relations.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger les relations CMDB";
  }
};

const fetchApplications = async () => {
  try {
    const response = await api.get("/cmdb/applications", headers());
    applications.value = response.data;
  } catch (err) {}
};

const fetchDatabases = async () => {
  try {
    const response = await api.get("/cmdb/databases", headers());
    databases.value = response.data;
  } catch (err) {}
};

const fetchAssets = async () => {
  try {
    const response = await api.get("/assets", headers());
    assets.value = response.data;
  } catch (err) {}
};

const getItemsByType = (type) => {
  switch (type) {
    case "Application": return applications.value;
    case "Asset": return assets.value;
    case "Database": return databases.value;
    default: return [];
  }
};

const getItemLabel = (type, id) => {
  const item = getItemsByType(type).find((el) => Number(el.ID) === Number(id));
  if (!item) return `${type} #${id}`;
  return item.name || `${type} #${id}`;
};

const resetForm = () => {
  form.value = { source_type: "Application", source_id: "", relation_type: "hosted_on", target_type: "Asset", target_id: "" };
  isEditMode.value = false;
  editingRelationId.value = null;
  error.value = "";
};

const openCreateModal = () => { resetForm(); showModal.value = true; };

const openEditModal = (relation) => {
  isEditMode.value = true;
  editingRelationId.value = relation.ID;
  form.value = {
    source_type: relation.source_type,
    source_id: relation.source_id,
    relation_type: relation.relation_type,
    target_type: relation.target_type,
    target_id: relation.target_id,
  };
  showModal.value = true;
};

const saveRelation = async () => {
  try {
    isSubmitting.value = true;
    error.value = "";
    const payload = {
      ...form.value,
      source_id: Number(form.value.source_id),
      target_id: Number(form.value.target_id),
    };

    if (isEditMode.value) {
      await api.put(`/cmdb/relations/${editingRelationId.value}`, payload, headers());
    } else {
      await api.post("/cmdb/relations", payload, headers());
    }
    showModal.value = false;
    resetForm();
    await fetchRelations();
  } catch (err) {
    error.value = "Impossible de valider le couplage des composants";
  } finally {
    isSubmitting.value = false;
  }
};

const deleteRelation = async (relation) => {
  if (!confirm("Supprimer cette relation du graphe de dépendance ?")) return;
  try {
    await api.delete(`/cmdb/relations/${relation.ID}`, headers());
    await fetchRelations();
  } catch (err) {
    error.value = "Échec de l'altération de la topologie";
  }
};

const filteredRelations = computed(() => {
  if (!searchQuery.value.trim()) return relations.value;
  const q = searchQuery.value.toLowerCase();

  return relations.value.filter((rel) => {
    const source = getItemLabel(rel.source_type, rel.source_id).toLowerCase();
    const target = getItemLabel(rel.target_type, rel.target_id).toLowerCase();

    return (
        source.includes(q) ||
        target.includes(q) ||
        rel.relation_type?.toLowerCase().includes(q) ||
        rel.source_type?.toLowerCase().includes(q) ||
        rel.target_type?.toLowerCase().includes(q)
    );
  });
});

const stats = computed(() => ({
  total: relations.value.length,
  hosted: relations.value.filter((r) => r.relation_type === "hosted_on").length,
  database: relations.value.filter((r) => r.relation_type === "uses_database").length,
  depends: relations.value.filter((r) => r.relation_type === "depends_on").length,
}));

onMounted(() => {
  fetchApplications();
  fetchDatabases();
  fetchAssets();
  fetchRelations();
});
</script>

<template>
  <DashboardLayout>
    <div class="topbar flex items-center justify-between">
      <div>
        <span class="tb-env">CMDB / Graphe de dépendances</span>
        <h1>Relations</h1>
        <p class="text-sm text-slate-500 mt-1">
          Cartographiez les dépendances structurelles entre vos applications, assets et bases de données.
        </p>
      </div>

      <button v-if="role === 'Admin'" @click="openCreateModal" class="btn-primary flex items-center gap-2">
        <span>+</span> Ajouter une relation
      </button>
    </div>

    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-muted flex items-center justify-center text-base w-11 h-11">Σ</div>
        <div>
          <div class="kc-lbl text-xs">Total Relations</div>
          <div class="kc-val text-xl font-bold mt-0.5">{{ stats.total }}</div>
        </div>
      </div>
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-brand flex items-center justify-center text-base w-11 h-11">☁</div>
        <div>
          <div class="kc-lbl text-xs">Hébergements</div>
          <div class="kc-val text-xl font-bold mt-0.5" style="color: var(--brand-light)">{{ stats.hosted }}</div>
        </div>
      </div>
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-info flex items-center justify-center text-base w-11 h-11">⛃</div>
        <div>
          <div class="kc-lbl text-xs">Bases liées</div>
          <div class="kc-val text-xl font-bold mt-0.5" style="color: var(--c-info-t)">{{ stats.database }}</div>
        </div>
      </div>
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-warn flex items-center justify-center text-base w-11 h-11">⚡</div>
        <div>
          <div class="kc-lbl text-xs">Dépendances</div>
          <div class="kc-val text-xl font-bold mt-0.5" style="color: var(--c-warn-t)">{{ stats.depends }}</div>
        </div>
      </div>
    </div>

    <div class="fbar flex items-center relative mb-6">
      <input
          v-model="searchQuery"
          type="text"
          placeholder="Rechercher une source, une cible ou un type de relation (ex: hosted_on, asset...)"
          class="search-input w-full pl-4 pr-10 py-2.5 text-sm"
      />
      <div class="absolute right-4 top-3 text-sm pointer-events-none">🔍</div>
    </div>

    <div v-if="error" class="modal-error error-bar mb-6 flex items-center gap-2">
      <span>⚠️</span> {{ error }}
    </div>

    <div class="table-card overflow-hidden cmdb-dt-wrapper">
      <DataTable
          :data="filteredRelations"
          :columns="columns"
          class="w-full text-left border-collapse"
          :options="dtOptions"
      >
        <template #column-source_id="{ rowData }">
          <div class="font-medium" style="color: var(--tx-primary);">
            {{ getItemLabel(rowData.source_type, rowData.source_id) }}
          </div>
          <div class="text-[11px]" style="color: var(--tx-muted);">
            {{ rowData.source_type }}
          </div>
        </template>

        <template #column-relation_type="{ cellData }">
          <span :class="['pill text-xs font-mono px-2.5 py-0.5 rounded-full border', relationConfig[cellData]?.class || 'p-sl']">
            {{ relationConfig[cellData]?.label || cellData }}
          </span>
        </template>

        <template #column-target_id="{ rowData }">
          <div class="font-medium" style="color: var(--tx-primary);">
            {{ getItemLabel(rowData.target_type, rowData.target_id) }}
          </div>
          <div class="text-[11px]" style="color: var(--tx-muted);">
            {{ rowData.target_type }}
          </div>
        </template>

        <template #column-3="{ rowData }">
          <div class="flex justify-end gap-1.5">
            <button
                v-if="role === 'Admin'"
                class="action-btn action-edit text-xs"
                title="Ajuster la liaison"
                @click="openEditModal(rowData)"
            >✏️</button>
            <button
                v-if="role === 'Admin'"
                class="action-btn action-delete text-xs"
                title="Supprimer la dépendance"
                @click="deleteRelation(rowData)"
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
            <span class="modal-eyebrow text-xs font-bold uppercase tracking-wider block mb-0.5" style="color: var(--brand-light)">Topologie Réseau</span>
            <h2 class="modal-title text-lg font-semibold m-0">
              {{ isEditMode ? "Modifier la relation" : "Ajouter une relation" }}
            </h2>
          </div>
          <button @click="showModal = false" class="modal-close p-2">✕</button>
        </div>

        <form @submit.prevent="saveRelation">
          <div class="p-5 space-y-5 overflow-y-auto max-h-[70vh]">

            <div class="p-4 rounded-xl border" style="background: var(--bg-base); border-color: var(--bd-subtle);">
              <span class="block text-xs font-bold uppercase tracking-wider mb-3" style="color: var(--tx-muted)">Composant Source (Départ)</span>
              <div class="grid grid-cols-2 gap-4">
                <div class="flex flex-col gap-1.5">
                  <label class="text-xs font-medium">Type CI</label>
                  <select v-model="form.source_type" class="field-input w-full appearance-none cursor-pointer" @change="form.source_id = ''">
                    <option v-for="type in typeOptions" :key="type" :value="type">{{ type }}</option>
                  </select>
                </div>
                <div class="flex flex-col gap-1.5">
                  <label class="text-xs font-medium">Instance</label>
                  <select v-model="form.source_id" required class="field-input w-full appearance-none cursor-pointer">
                    <option value="" disabled selected>Sélectionner la source</option>
                    <option v-for="item in getItemsByType(form.source_type)" :key="item.ID" :value="item.ID">
                      {{ item.name }}
                    </option>
                  </select>
                </div>
              </div>
            </div>

            <div class="flex flex-col gap-1.5 px-1">
              <label class="text-xs font-medium">Nature de l'interdépendance (Verbe de liaison)</label>
              <select v-model="form.relation_type" required class="field-input w-full appearance-none cursor-pointer">
                <option v-for="relation in relationOptions" :key="relation" :value="relation">{{ relation }}</option>
              </select>
            </div>

            <div class="p-4 rounded-xl border" style="background: var(--bg-base); border-color: var(--bd-subtle);">
              <span class="block text-xs font-bold uppercase tracking-wider mb-3" style="color: var(--tx-muted)">Composant Cible (Arrivée)</span>
              <div class="grid grid-cols-2 gap-4">
                <div class="flex flex-col gap-1.5">
                  <label class="text-xs font-medium">Type CI</label>
                  <select v-model="form.target_type" class="field-input w-full appearance-none cursor-pointer" @change="form.target_id = ''">
                    <option v-for="type in typeOptions" :key="type" :value="type">{{ type }}</option>
                  </select>
                </div>
                <div class="flex flex-col gap-1.5">
                  <label class="text-xs font-medium">Instance</label>
                  <select v-model="form.target_id" required class="field-input w-full appearance-none cursor-pointer">
                    <option value="" disabled selected>Sélectionner la cible</option>
                    <option v-for="item in getItemsByType(form.target_type)" :key="item.ID" :value="item.ID">
                      {{ item.name }}
                    </option>
                  </select>
                </div>
              </div>
            </div>

          </div>

          <div class="modal-footer p-4 flex justify-end gap-3 border-t" style="border-color: var(--bd-subtle); background: rgba(0,0,0,0.05)">
            <button type="button" @click="showModal = false" class="btn-ghost">Annuler</button>
            <button type="submit" :disabled="isSubmitting" class="btn-primary min-w-[100px]">
              <span v-if="isSubmitting" class="inline-block spinner w-4 h-4 rounded-full border-2 animate-spin mr-2" style="vertical-align: sub;"></span>
              {{ isEditMode ? "Enregistrer" : "Créer le lien" }}
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

.spinner {
  border-color: transparent;
  border-top-color: #fff;
}
</style>

<style>
/* Encapsulation globale DataTables mappée sur ton UI dark system */
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
