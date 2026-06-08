<script setup>
import { onMounted, ref, computed } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

import DataTable from "datatables.net-vue3";
import DataTablesCore from "datatables.net-dt";

import "datatables.net-dt/css/dataTables.dataTables.css";

DataTable.use(DataTablesCore);

const sites = ref([]);
const showModal = ref(false);
const isEditMode = ref(false);
const editingSiteId = ref(null);
const error = ref("");
const isSubmitting = ref(false);
const searchQuery = ref("");

const form = ref({
  name: "",
  city: "",
  country: "France",
  description: "",
});

const getToken = () => localStorage.getItem("token");

// Configuration des colonnes DataTable (Immunisées contre le Warning n°4 de l'API)
const columns = [
  { data: "name", title: "Nom", defaultContent: "Sans nom" },
  { data: "city", title: "Ville", defaultContent: "—" },
  { data: "country", title: "Pays", defaultContent: "" },
  { data: "description", title: "Description", defaultContent: "—" },
  { data: null, title: "Actions", orderable: false, className: "text-right" }
];

const dtOptions = {
  responsive: true,
  pageLength: 10,
  dom: 't<"dt-footer flex items-center justify-between p-4"ip>', // Recherche native masquée pour utiliser la tienne
  language: {
    info: "Affichage de _START_ à _END_ sur _TOTAL_ implantations géographiques",
    infoEmpty: "Aucun site enregistré",
    infoFiltered: "(filtré depuis _MAX_ éléments)",
    zeroRecords: "Aucun site trouvé",
    paginate: {
      next: "➔",
      previous: " Zar "
    }
  }
};

const fetchSites = async () => {
  try {
    const response = await api.get("/sites", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    sites.value = response.data || [];
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger le registre des sites physiques";
  }
};

// Filtrage réactif lié à ta barre de recherche personnalisée
const filteredSites = computed(() => {
  let list = sites.value;
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    list = list.filter(
        (site) =>
            site.name?.toLowerCase().includes(q) ||
            site.city?.toLowerCase().includes(q) ||
            site.country?.toLowerCase().includes(q) ||
            site.description?.toLowerCase().includes(q)
    );
  }
  return list;
});

const resetForm = () => {
  form.value = { name: "", city: "", country: "France", description: "" };
  isEditMode.value = false;
  editingSiteId.value = null;
  error.value = "";
};

const openCreateModal = () => { resetForm(); showModal.value = true; };

const openEditModal = (site) => {
  isEditMode.value = true;
  editingSiteId.value = site.ID;
  form.value = {
    name: site.name || "",
    city: site.city || "",
    country: site.country || "France",
    description: site.description || "",
  };
  showModal.value = true;
};

const saveSite = async () => {
  try {
    isSubmitting.value = true;
    error.value = "";

    if (isEditMode.value) {
      await api.put(`/sites/${editingSiteId.value}`, form.value, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    } else {
      await api.post("/sites", form.value, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    }

    showModal.value = false;
    resetForm();
    await fetchSites();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de sauvegarder les modifications du site";
  } finally {
    isSubmitting.value = false;
  }
};

const deleteSite = async (site) => {
  if (!confirm(`Confirmez-vous la suppression définitive du site "${site.name}" ?`)) return;

  try {
    await api.delete(`/sites/${site.ID}`, {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    await fetchSites();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de purger ce site (Vérifiez les liaisons ou dépendances actives)";
  }
};

onMounted(fetchSites);
</script>

<template>
  <DashboardLayout>
    <div class="topbar flex items-center justify-between">
      <div>
        <span class="tb-env">Infrastructure / Datacenters & Bureaux</span>
        <h1>Sites</h1>
        <p class="text-sm text-slate-500 mt-1">Gérez les implantations géographiques et infrastructures de votre organisation.</p>
      </div>

      <button @click="openCreateModal" class="btn-primary flex items-center gap-2">
        <span>+</span> Ajouter un site
      </button>
    </div>

    <div class="fbar flex items-center justify-between gap-4 flex-wrap mb-6">
      <div class="search-wrap flex-1 relative">
        <input
            v-model="searchQuery"
            type="text"
            placeholder="Rechercher par nom de site, ville, pays ou description..."
            class="search-input w-full pl-4 pr-10 py-2.5"
        />
        <span class="absolute right-4 top-3 text-sm">🔍</span>
      </div>
    </div>

    <div v-if="error" class="modal-error error-bar mb-6 flex items-center gap-2">
      <span>⚠️</span> {{ error }}
    </div>

    <div class="table-card overflow-hidden cmdb-dt-wrapper">
      <DataTable
          :data="filteredSites"
          :columns="columns"
          class="w-full text-left border-collapse"
          :options="dtOptions"
      >
        <template #column-name="{ cellData }">
          <div class="font-medium" style="color: var(--tx-primary);">{{ cellData }}</div>
        </template>

        <template #column-city="{ cellData }">
          <span class="font-medium" style="color: var(--tx-secondary);">{{ cellData }}</span>
        </template>

        <template #column-country="{ cellData }">
          <span v-if="cellData" class="sc-bdg bdg-sl px-2 py-0.5 rounded text-xs font-medium border">
            {{ cellData }}
          </span>
          <span v-else class="text-slate-600">—</span>
        </template>

        <template #column-description="{ cellData }">
          <span class="text-sm text-slate-400 line-clamp-1 max-w-xs" :title="cellData">{{ cellData }}</span>
        </template>

        <template #column-4="{ rowData }">
          <div class="flex justify-end gap-1.5">
            <button
                class="action-btn action-edit text-xs"
                title="Modifier le site"
                @click="openEditModal(rowData)"
            >✏️</button>
            <button
                class="action-btn action-delete text-xs"
                title="Supprimer le site"
                @click="deleteSite(rowData)"
            >🗑️</button>
          </div>
        </template>
      </DataTable>
    </div>
  </DashboardLayout>

  <Teleport to="body">
    <div v-if="showModal" class="modal-overlay fixed inset-0 flex items-center justify-center p-4 z-50" @click.self="showModal = false">
      <div class="modal rounded-2xl w-full max-w-lg overflow-hidden shadow-2xl">

        <div class="modal-header p-5 flex justify-between items-center border-b" style="border-color: var(--bd-subtle)">
          <div>
            <span class="modal-eyebrow text-xs font-bold uppercase tracking-wider block mb-0.5" style="color: var(--brand-light)">Infrastructures</span>
            <h2 class="modal-title text-lg font-semibold m-0">
              {{ isEditMode ? "Modifier le site physique" : "Ajouter une implantation" }}
            </h2>
          </div>
          <button @click="showModal = false" class="modal-close p-2">✕</button>
        </div>

        <form @submit.prevent="saveSite">
          <div class="p-5 space-y-4 overflow-y-auto max-h-[70vh]">

            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium">Nom du site <span style="color: #f43f5e;">*</span></label>
              <input
                  v-model="form.name"
                  required
                  placeholder="Ex: Siège Social, Datacenter DC2..."
                  class="field-input w-full"
              />
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-medium">Ville <span style="color: #f43f5e;">*</span></label>
                <input
                    v-model="form.city"
                    required
                    placeholder="Ex: Paris"
                    class="field-input w-full"
                />
              </div>
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-medium">Pays <span style="color: #f43f5e;">*</span></label>
                <input
                    v-model="form.country"
                    required
                    placeholder="Ex: France"
                    class="field-input w-full"
                />
              </div>
            </div>

            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium">Description de la plaque / Infos d'adressage</label>
              <textarea
                  v-model="form.description"
                  placeholder="Coordonnées d'accès, armoires réseau ou détails opérationnels..."
                  rows="3"
                  class="field-input w-full resize-none"
              ></textarea>
            </div>

          </div>

          <div class="modal-footer p-4 flex justify-end gap-3 border-t" style="border-color: var(--bd-subtle); background: rgba(0,0,0,0.05)">
            <button type="button" @click="showModal = false" class="btn-ghost">
              Annuler
            </button>
            <button type="submit" :disabled="isSubmitting" class="btn-primary min-w-[130px]">
              <span v-if="isSubmitting" class="inline-block spinner w-4 h-4 rounded-full border-2 animate-spin mr-2" style="vertical-align: sub;"></span>
              {{ isEditMode ? "Enregistrer" : "Créer le site" }}
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
/* Alignement global DataTables sur la charte graphique globale */
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