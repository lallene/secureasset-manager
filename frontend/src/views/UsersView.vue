<script setup>
import { onMounted, ref, computed } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

import DataTable from "datatables.net-vue3";
import DataTablesCore from "datatables.net-dt";

import "datatables.net-dt/css/dataTables.dataTables.css";

DataTable.use(DataTablesCore);

const applications = ref([]);
const sites = ref([]);

const showModal = ref(false);
const isEditMode = ref(false);
const editingApplicationId = ref(null);

const error = ref("");
const isSubmitting = ref(false);
const searchQuery = ref("");
const filterCriticity = ref("all");

const form = ref({
  name: "",
  code: "",
  criticity: "P3",
  environment: "Production",
  site_id: "",
});

const getToken = () => localStorage.getItem("token");
const role = computed(() => localStorage.getItem("role") || "");

const criticityOptions = ["P1", "P2", "P3", "P4"];
const environmentOptions = ["Production", "Preproduction", "Staging", "Development", "Test"];

// Configuration visuelle des criticités CMDB (SLA / Impact)
const criticityConfig = {
  P1: { class: "bdg-rose pulse-slow font-bold", label: "P1 - Critique" },
  P2: { class: "bdg-amb font-semibold",        label: "P2 - Élevé" },
  P3: { class: "bdg-ind",                      label: "P3 - Moyen" },
  P4: { class: "bdg-sl",                       label: "P4 - Faible" }
};

const environmentConfig = {
  Production:    { class: "p-rose", label: "Production" },
  Preproduction: { class: "p-ind",  label: "Preprod" },
  Staging:       { class: "p-amb",  label: "Staging" },
  Development:   { class: "p-sl",   label: "Dev" },
  Test:          { class: "p-sl",   label: "Test" }
};

const columns = [
  { data: "name", title: "Application", defaultContent: "Sans nom" },
  { data: "code", title: "Code Trigramme", defaultContent: "" },
  { data: "criticity", title: "Criticité (SLA)", defaultContent: "P3" }, // ◄ Sécurisé ici avec "P3" par défaut
  { data: "environment", title: "Environnement", defaultContent: "Production" }, // ◄ Sécurisé ici
  { data: "site.name", title: "Site Principal", defaultContent: "—" },
  { data: null, title: "Actions", orderable: false, className: "text-right" }
];

const dtOptions = {
  responsive: true,
  pageLength: 10,
  dom: 't<"dt-footer flex items-center justify-between p-4"ip>',
  language: {
    info: "Affichage de _START_ à _END_ sur _TOTAL_ applications cataloguées",
    infoEmpty: "Aucune application répertoriée",
    infoFiltered: "(filtré depuis _MAX_ éléments)",
    zeroRecords: "Aucun asset applicatif trouvé",
    paginate: {
      next: "➔",
      previous: " Zar "
    }
  }
};

const fetchApplications = async () => {
  try {
    const response = await api.get("/cmdb/applications", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    applications.value = response.data || [];
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger le catalogue des applications";
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
  }
};

const filteredApplications = computed(() => {
  let list = applications.value;

  if (filterCriticity.value !== "all") {
    list = list.filter((app) => app.criticity === filterCriticity.value);
  }

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    list = list.filter(
        (app) =>
            app.name?.toLowerCase().includes(q) ||
            app.code?.toLowerCase().includes(q) ||
            app.criticity?.toLowerCase().includes(q) ||
            app.environment?.toLowerCase().includes(q) ||
            app.site?.name?.toLowerCase().includes(q)
    );
  }

  return list;
});

const stats = computed(() => ({
  total: applications.value.length,
  p1: applications.value.filter((app) => app.criticity === "P1").length,
  production: applications.value.filter((app) => app.environment === "Production").length,
}));

const resetForm = () => {
  form.value = { name: "", code: "", criticity: "P3", environment: "Production", site_id: "" };
  isEditMode.value = false;
  editingApplicationId.value = null;
  error.value = "";
};

const openCreateModal = () => { resetForm(); showModal.value = true; };

const openEditModal = (app) => {
  isEditMode.value = true;
  editingApplicationId.value = app.ID;
  form.value = {
    name: app.name || "",
    code: app.code || "",
    criticity: app.criticity || "P3",
    environment: app.environment || "Production",
    site_id: app.site_id || "",
  };
  showModal.value = true;
};

const saveApplication = async () => {
  try {
    isSubmitting.value = true;
    error.value = "";
    const payload = {
      ...form.value,
      site_id: form.value.site_id ? Number(form.value.site_id) : null,
    };

    if (isEditMode.value) {
      await api.put(`/cmdb/applications/${editingApplicationId.value}`, payload, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    } else {
      await api.post("/cmdb/applications", payload, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    }
    showModal.value = false;
    resetForm();
    await fetchApplications();
  } catch (err) {
    console.error(err);
    error.value = "Erreur lors de la mise à jour du catalogue applicatif";
  } finally {
    isSubmitting.value = false;
  }
};

const deleteApplication = async (app) => {
  if (!confirm(`Supprimer l'application "${app.name}" ? (Attention aux liaisons CMDB actives)`)) return;
  try {
    await api.delete(`/cmdb/applications/${app.ID}`, {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    await fetchApplications();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de purger cet élément (Contraintes d'intégrité en cours)";
  }
};

onMounted(() => {
  fetchApplications();
  fetchSites();
});
</script>

<template>
  <DashboardLayout>
    <div class="topbar flex items-center justify-between">
      <div>
        <span class="tb-env">CMDB / Software Asset Management</span>
        <h1>Applications</h1>
        <p class="text-sm text-slate-500 mt-1">Inventoriez et qualifiez le niveau de criticité des applications métiers.</p>
      </div>

      <button v-if="role === 'Admin'" @click="openCreateModal" class="btn-primary flex items-center gap-2">
        <span>+</span> Ajouter une application
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-muted flex items-center justify-center text-base w-11 h-11">💻</div>
        <div>
          <div class="kc-lbl text-xs">Applications Recensées</div>
          <div class="kc-val text-xl font-bold mt-0.5">{{ stats.total }}</div>
        </div>
      </div>
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-danger flex items-center justify-center text-base w-11 h-11">🔥</div>
        <div>
          <div class="kc-lbl text-xs">Services Critiques (P1)</div>
          <div class="kc-val text-xl font-bold mt-0.5" style="color: #f43f5e;">{{ stats.p1 }}</div>
        </div>
      </div>
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-brand flex items-center justify-center text-base w-11 h-11">🌐</div>
        <div>
          <div class="kc-lbl text-xs">Déployées en Production</div>
          <div class="kc-val text-xl font-bold mt-0.5" style="color: var(--brand-light)">{{ stats.production }}</div>
        </div>
      </div>
    </div>

    <div class="fbar flex items-center justify-between gap-4 flex-wrap mb-6">
      <div class="search-wrap flex-1 relative">
        <input
            v-model="searchQuery"
            type="text"
            placeholder="Rechercher par nom d'application, trigramme code, site d'hébergement..."
            class="search-input w-full pl-4 pr-10 py-2.5"
        />
        <span class="absolute right-4 top-3 text-sm">🔍</span>
      </div>

      <select v-model="filterCriticity" class="field-input w-52">
        <option value="all">Toutes les criticités</option>
        <option v-for="crit in criticityOptions" :key="crit" :value="crit">Criticité {{ crit }}</option>
      </select>
    </div>

    <div v-if="error" class="modal-error error-bar mb-6 flex items-center gap-2">
      <span>⚠️</span> {{ error }}
    </div>

    <div class="table-card overflow-hidden cmdb-dt-wrapper">
      <DataTable
          :data="filteredApplications"
          :columns="columns"
          class="w-full text-left border-collapse"
          :options="dtOptions"
      >
        <template #column-name="{ cellData }">
          <div class="font-medium" style="color: var(--tx-primary)">{{ cellData }}</div>
        </template>

        <template #column-code="{ cellData }">
          <span class="px-2 py-0.5 font-mono text-xs rounded border" style="background: var(--bg-base); border-color: var(--bd-subtle); color: var(--tx-muted)">
            {{ cellData || '—' }}
          </span>
        </template>

        <template #column-criticity="{ cellData }">
          <span :class="['sc-bdg text-xs px-2.5 py-0.5 rounded-full border', criticityConfig[cellData]?.class || 'bdg-sl']">
            {{ criticityConfig[cellData]?.label || cellData }}
          </span>
        </template>

        <template #column-environment="{ cellData }">
          <span :class="['pill text-xs font-medium px-2.5 py-0.5 rounded-full border', environmentConfig[cellData]?.class || 'p-sl']">
            {{ environmentConfig[cellData]?.label || cellData }}
          </span>
        </template>

        <template #column-5="{ rowData }">
          <div class="flex justify-end gap-1.5">
            <button
                v-if="role === 'Admin'"
                class="action-btn action-edit text-xs"
                title="Ajuster l'application"
                @click="openEditModal(rowData)"
            >✏️</button>
            <button
                v-if="role === 'Admin'"
                class="action-btn action-delete text-xs"
                title="Retirer du catalogue informatique"
                @click="deleteApplication(rowData)"
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
            <span class="modal-eyebrow text-xs font-bold uppercase tracking-wider block mb-0.5" style="color: var(--brand-light)">Registre Logiciel</span>
            <h2 class="modal-title text-lg font-semibold m-0">
              {{ isEditMode ? "Modifier la fiche applicative" : "Enregistrer une nouvelle application" }}
            </h2>
          </div>
          <button @click="showModal = false" class="modal-close p-2">✕</button>
        </div>

        <form @submit.prevent="saveApplication">
          <div class="p-5 space-y-4 overflow-y-auto max-h-[70vh]">

            <div class="grid grid-cols-3 gap-4">
              <div class="col-span-2 flex flex-col gap-1.5">
                <label class="text-xs font-medium">Nom de l'application <span style="color: #f43f5e;">*</span></label>
                <input v-model="form.name" required type="text" class="field-input w-full" placeholder="Ex: Active Directory Auth" />
              </div>

              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-medium">Code (Trigramme) <span style="color: #f43f5e;">*</span></label>
                <input v-model="form.code" required type="text" class="field-input w-full font-mono uppercase" placeholder="Ex: ADA" maxLength="5" />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-medium">Seuil de criticité SLA <span style="color: #f43f5e;">*</span></label>
                <select v-model="form.criticity" required class="field-input w-full appearance-none cursor-pointer">
                  <option value="P1">P1 - Interruption Critique (Bloquant)</option>
                  <option value="P2">P2 - Incident Majeur (Dégradé)</option>
                  <option value="P3">P3 - Impact Moyen (Standard)</option>
                  <option value="P4">P4 - Demande Mineure (Consultatif)</option>
                </select>
              </div>

              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-medium">Scope d'environnement <span style="color: #f43f5e;">*</span></label>
                <select v-model="form.environment" required class="field-input w-full appearance-none cursor-pointer">
                  <option v-for="env in environmentOptions" :key="env" :value="env">{{ env }}</option>
                </select>
              </div>
            </div>

            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium">Site ou Cluster d'hébergement principal</label>
              <select v-model="form.site_id" class="field-input w-full appearance-none cursor-pointer">
                <option value="">Sélectionner un site physique</option>
                <option v-for="site in sites" :key="site.ID" :value="site.ID">{{ site.name }}</option>
              </select>
            </div>

          </div>

          <div class="modal-footer p-4 flex justify-end gap-3 border-t" style="border-color: var(--bd-subtle); background: rgba(0,0,0,0.05)">
            <button type="button" class="btn-ghost" @click="showModal = false">Annuler</button>
            <button type="submit" :disabled="isSubmitting" class="btn-primary min-w-[120px]">
              <span v-if="isSubmitting" class="inline-block spinner w-4 h-4 rounded-full border-2 animate-spin mr-2" style="vertical-align: sub;"></span>
              {{ isEditMode ? "Mettre à jour" : "Déployer la fiche" }}
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

.spinner {
  border-color: transparent;
  border-top-color: #fff;
}

@keyframes pulseGlow {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}
</style>

<style>
/* Encapsulation globale DataTables partagée sur l'UI system dark */
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