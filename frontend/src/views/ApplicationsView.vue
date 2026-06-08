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
  const users = ref([]);

  const showModal = ref(false);
  const isEditMode = ref(false);
  const editingApplicationId = ref(null);

  const error = ref("");
  const isSubmitting = ref(false);
  const searchQuery = ref("");
  const filterCriticality = ref("all");

  const form = ref({
    name: "",
    description: "",
    version: "",
    criticality: "Medium",
    status: "Active",
    site_id: "",
    owner_id: "",
  });

  const getToken = () => localStorage.getItem("token");
  const role = computed(() => localStorage.getItem("role") || "");

  const criticalityOptions = ["Low", "Medium", "High", "Critical"];
  const statusOptions = ["Active", "Maintenance", "Retired"];

  // Mapping sémantique aligné sur tes classes de badges globales
  const criticalityConfig = {
    Low:      { class: "p-grn", label: "Faible" },
    Medium:   { class: "p-ind", label: "Moyenne" },
    High:     { class: "p-amb", label: "Haute" },
    Critical: { class: "p-rose pulse-slow", label: "Critique" },
  };

  const statusConfig = {
    Active:      { class: "bdg-brand", label: "Actif" },
    Maintenance: { class: "bdg-ind",   label: "Maintenance" },
    Retired:     { class: "bdg-sl",    label: "Retiré" },
  };

  const columns = [
    { data: "name", title: "Application" },
    { data: "version", title: "Version" },
    { data: "criticality", title: "Criticité" },
    { data: "status", title: "Statut" },
    { data: "site.name", title: "Site affecté", defaultContent: "—" },
    { data: "owner.name", title: "Responsable (Owner)", defaultContent: "—" },
    { data: null, title: "Actions", orderable: false, className: "text-right" }
  ];

  // Options avancées DataTables injectées pour l'override du DOM natif
  const dtOptions = {
    responsive: true,
    pageLength: 10,
    dom: 't<"dt-footer flex items-center justify-between p-4"ip>',
    language: {
      info: "Affichage de _START_ à _END_ sur _TOTAL_ architectures",
      infoEmpty: "Aucune application enregistrée",
      infoFiltered: "(filtré depuis _MAX_ éléments)",
      zeroRecords: "Aucune correspondance trouvée",
      paginate: {
        next: "➔",
        previous: " Zar "
      }
    }
  };

  const filteredApplications = computed(() => {
    let list = applications.value;
    if (filterCriticality.value !== "all") {
      list = list.filter((app) => app.criticality === filterCriticality.value);
    }
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase();
      list = list.filter(
          (app) =>
              app.name?.toLowerCase().includes(q) ||
              app.version?.toLowerCase().includes(q) ||
              app.site?.name?.toLowerCase().includes(q) ||
              app.owner?.name?.toLowerCase().includes(q)
      );
    }
    return list;
  });

  const stats = computed(() => ({
    total: applications.value.length,
    critical: applications.value.filter((a) => a.criticality === "Critical").length,
    high: applications.value.filter((a) => a.criticality === "High").length,
    active: applications.value.filter((a) => a.status === "Active").length,
  }));

  const fetchApplications = async () => {
    try {
      const response = await api.get("/cmdb/applications", {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
      applications.value = response.data;
    } catch (err) {
      console.error(err);
      error.value = "Impossible de charger la cartographie applicative";
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

  const fetchUsers = async () => {
    try {
      const response = await api.get("/users", {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
      users.value = response.data;
    } catch (err) {}
  };

  const resetForm = () => {
    form.value = { name: "", description: "", version: "", criticality: "Medium", status: "Active", site_id: "", owner_id: "" };
    isEditMode.value = false;
    editingApplicationId.value = null;
    error.value = "";
  };

  const openCreateModal = () => { resetForm(); showModal.value = true; };

  const openEditModal = (app) => {
    isEditMode.value = true;
    editingApplicationId.value = app.ID;
    form.value = {
      name: app.name,
      description: app.description,
      version: app.version,
      criticality: app.criticality,
      status: app.status,
      site_id: app.site_id,
      owner_id: app.owner_id || "",
    };
    showModal.value = true;
  };

  const saveApplication = async () => {
    try {
      isSubmitting.value = true;
      error.value = "";
      const payload = {
        ...form.value,
        site_id: Number(form.value.site_id),
        owner_id: form.value.owner_id ? Number(form.value.owner_id) : null,
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
      error.value = "Erreur lors de la sauvegarde de la fiche applicative";
    } finally {
      isSubmitting.value = false;
    }
  };

  const deleteApplication = async (app) => {
    if (!confirm(`Supprimer l'application "${app.name}" du catalogue CMDB ?`)) return;
    try {
      await api.delete(`/cmdb/applications/${app.ID}`, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
      await fetchApplications();
    } catch (err) {
      error.value = "Rupture d'intégrité : impossible de purger cet item";
    }
  };

  onMounted(() => {
    fetchApplications();
    fetchSites();
    fetchUsers();
  });
</script>

<template>
  <DashboardLayout>
    <div class="topbar flex items-center justify-between">
      <div>
        <span class="tb-env">Configuration Management Database (CMDB)</span>
        <h1>Applications</h1>
        <p class="text-sm text-slate-500 mt-1">
          Cartographiez les dépendances, versions et niveaux de criticité de votre écosystème logiciel.
        </p>
      </div>

      <button
          v-if="role === 'Admin'"
          @click="openCreateModal"
          class="btn-primary flex items-center gap-2"
      >
        <span>+</span> Ajouter une application
      </button>
    </div>

    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-muted flex items-center justify-center text-base w-11 h-11">📦</div>
        <div>
          <div class="kc-lbl text-xs">Total Composants</div>
          <div class="kc-val text-xl font-bold mt-0.5">{{ stats.total }}</div>
        </div>
      </div>
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-brand flex items-center justify-center text-base w-11 h-11">💻</div>
        <div>
          <div class="kc-lbl text-xs">Actives</div>
          <div class="kc-val text-xl font-bold mt-0.5" style="color: var(--brand-light)">{{ stats.active }}</div>
        </div>
      </div>
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-warn flex items-center justify-center text-base w-11 h-11">⚠️</div>
        <div>
          <div class="kc-lbl text-xs">Criticité Haute</div>
          <div class="kc-val text-xl font-bold mt-0.5" style="color: var(--c-warn-t)">{{ stats.high }}</div>
        </div>
      </div>
      <div class="kc flex items-center gap-4">
        <div class="kc-ico ico-danger flex items-center justify-center text-base w-11 h-11">🚨</div>
        <div>
          <div class="kc-lbl text-xs">Périmètres Critiques</div>
          <div class="kc-val text-xl font-bold mt-0.5" style="color: var(--c-rose-t)">{{ stats.critical }}</div>
        </div>
      </div>
    </div>

    <div class="fbar flex items-center justify-between gap-4 flex-wrap mb-6">
      <div class="search-wrap flex-1 relative">
        <input
            v-model="searchQuery"
            type="text"
            placeholder="Rechercher par nom d'application, tag de version, site physique, owner référent..."
            class="search-input w-full pl-4 pr-10 py-2.5"
        />
        <span class="absolute right-4 top-3 text-sm">🔍</span>
      </div>

      <select v-model="filterCriticality" class="field-input w-52">
        <option value="all">Toutes criticités</option>
        <option v-for="c in criticalityOptions" :key="c" :value="c">
          Niveau : {{ criticalityConfig[c]?.label || c }}
        </option>
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
        <template #column-name="{ cellData, rowData }">
          <div class="font-medium" style="color: var(--tx-primary)">{{ cellData }}</div>
          <div class="text-[11px] truncate max-w-[200px]" style="color: var(--tx-muted)" :title="rowData.description">
            {{ rowData.description || 'Aucune documentation technique' }}
          </div>
        </template>

        <template #column-version="{ cellData }">
          <span class="px-1.5 py-0.5 rounded text-[11px] font-mono border" style="background: var(--bg-base); border-color: var(--bd-subtle); color: var(--tx-secondary)">
            v{{ cellData || '0.0.0' }}
          </span>
        </template>

        <template #column-criticality="{ cellData }">
          <span :class="['pill text-xs font-medium px-2.5 py-0.5 rounded-full border', criticalityConfig[cellData]?.class || 'p-sl']">
            {{ criticalityConfig[cellData]?.label || cellData }}
          </span>
        </template>

        <template #column-status="{ cellData }">
          <span :class="['sc-bdg text-xs px-2.5 py-0.5 rounded-full border font-medium', statusConfig[cellData]?.class || 'bdg-sl']">
            {{ statusConfig[cellData]?.label || cellData }}
          </span>
        </template>

        <template #column-6="{ rowData }">
          <div class="flex justify-end gap-1.5">
            <button
                class="action-btn action-edit text-xs"
                title="Éditer la configuration"
                @click="openEditModal(rowData)"
            >✏️</button>
            <button
                v-if="role === 'Admin'"
                class="action-btn action-delete text-xs"
                title="Purger du référentiel"
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
            <span class="modal-eyebrow text-xs font-bold uppercase tracking-wider block mb-0.5" style="color: var(--brand-light)">Topologie SI</span>
            <h2 class="modal-title text-lg font-semibold m-0">
              {{ isEditMode ? "Modifier la fiche applicative" : "Déclarer un nouveau composant" }}
            </h2>
          </div>
          <button @click="showModal = false" class="modal-close p-2">✕</button>
        </div>

        <form @submit.prevent="saveApplication">
          <div class="p-5 space-y-4 overflow-y-auto max-h-[70vh]">

            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium">Désignation du composant <span style="color: var(--c-rose-t)">*</span></label>
              <input v-model="form.name" required type="text" class="field-input w-full" placeholder="Ex: Active Directory LDAP Synchro" />
            </div>

            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium">Spécifications & Périmètre d'usage</label>
              <textarea v-model="form.description" rows="3" class="field-input w-full" placeholder="Rôle fonctionnel, dépendances critiques amont/aval..."></textarea>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-medium">Version Release</label>
                <input v-model="form.version" type="text" class="field-input w-full font-mono" placeholder="Ex: v2.4.1-stable" />
              </div>

              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-medium">Seuil de Criticité <span style="color: var(--c-rose-t)">*</span></label>
                <select v-model="form.criticality" required class="field-input w-full appearance-none cursor-pointer">
                  <option value="Low">Faible (Low)</option>
                  <option value="Medium">Modéré (Medium)</option>
                  <option value="High">Élevé (High)</option>
                  <option value="Critical">Vital / Critique (Critical)</option>
                </select>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-medium">Cycle de vie</label>
                <select v-model="form.status" class="field-input w-full appearance-none cursor-pointer">
                  <option v-for="s in statusOptions" :key="s" :value="s">{{ s }}</option>
                </select>
              </div>

              <div class="flex flex-col gap-1.5">
                <label class="text-xs font-medium">Site d'Hébergement / Cluster <span style="color: var(--c-rose-t)">*</span></label>
                <select v-model="form.site_id" required class="field-input w-full appearance-none cursor-pointer">
                  <option value="" disabled>Sélectionner un datacenter</option>
                  <option v-for="site in sites" :key="site.ID" :value="site.ID">{{ site.name }}</option>
                </select>
              </div>
            </div>

            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-medium">Ingénieur Responsable App (Owner)</label>
              <select v-model="form.owner_id" class="field-input w-full appearance-none cursor-pointer">
                <option value="">Aucune attribution de responsabilité</option>
                <option v-for="user in users" :key="user.ID" :value="user.ID">
                  {{ user.name }} — [{{ user.role }}]
                </option>
              </select>
            </div>
          </div>

          <div class="modal-footer p-4 flex justify-end gap-3 border-t" style="border-color: var(--bd-subtle); background: rgba(0,0,0,0.05)">
            <button type="button" class="btn-ghost" @click="showModal = false">Annuler</button>
            <button type="submit" :disabled="isSubmitting" class="btn-primary">
              {{ isEditMode ? "Enregistrer les modifications" : "Inscrire au registre" }}
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
  /* Overrides des injections globales de DataTables pour correspondre au thème de l'application */
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