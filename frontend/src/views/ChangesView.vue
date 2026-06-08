<script setup>
import { ref, computed, onMounted } from "vue";
import api from "../api/axios.js";
import DashboardLayout from "../layouts/DashboardLayout.vue";

import DataTable from "datatables.net-vue3";
import DataTablesCore from "datatables.net-dt";
import Buttons from "datatables.net-buttons-dt";
import "datatables.net-dt/css/dataTables.dataTables.css";
import "datatables.net-buttons-dt/css/buttons.dataTables.css";

DataTable.use(DataTablesCore);
DataTable.use(Buttons);

import { useRouter } from "vue-router";

const router = useRouter();

const changes = ref([]);
const businessServices = ref([]);

const showModal = ref(false);
const isEditMode = ref(false);
const editingChangeId = ref(null);
const error = ref("");
const isLoading = ref(false);

const role = computed(() => localStorage.getItem("role") || "");
const getToken = () => localStorage.getItem("token");

const form = ref({
  title: "",
  description: "",
  type: "Normal",
  risk: "Medium",
  status: "Draft",
  rollback_plan: "",
  business_service_id: null,
});

const columns = [
  { data: "reference", title: "Référence", defaultContent: "—" },
  { data: "title", title: "Titre" },
  { data: "type", title: "Type" },
  { data: "risk", title: "Risque" },
  { data: "status", title: "Statut" },
  { data: "business_service.name", title: "Service métier", defaultContent: "—" },
  { data: null, title: "Actions", orderable: false },
];

const dtOptions = {
  pageLength: 10,
  order: [[0, "desc"]],
  dom: '<"dt-top"Bf>t<"dt-footer"ip>',
  buttons: [
    {
      extend: "csvHtml5",
      text: "Exporter CSV",
      className: "dt-btn",
    },
    {
      extend: "excelHtml5",
      text: "Exporter Excel",
      className: "dt-btn",
    },
  ],
  language: {
    search: "Rechercher :",
    info: "_START_ à _END_ sur _TOTAL_ changements",
    zeroRecords: "Aucun changement trouvé",
    paginate: {
      previous: "←",
      next: "→",
    },
  },
};

const fetchChanges = async () => {
  isLoading.value = true;
  try {
    const response = await api.get("/changes", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    changes.value = response.data || [];
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger les changements";
  } finally {
    isLoading.value = false;
  }
};

const fetchBusinessServices = async () => {
  try {
    const response = await api.get("/cmdb/business-services", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    businessServices.value = response.data || [];
  } catch (err) {
    console.error(err);
  }
};

const resetForm = () => {
  form.value = {
    title: "",
    description: "",
    type: "Normal",
    risk: "Medium",
    status: "Draft",
    rollback_plan: "",
    business_service_id: null,
  };
  isEditMode.value = false;
  editingChangeId.value = null;
  error.value = "";
};

const openCreateModal = () => {
  resetForm();
  showModal.value = true;
};

const openEditModal = (change) => {
  isEditMode.value = true;
  editingChangeId.value = change.ID;

  form.value = {
    title: change.title || "",
    description: change.description || "",
    type: change.type || "Normal",
    risk: change.risk || "Medium",
    status: change.status || "Draft",
    rollback_plan: change.rollback_plan || "",
    business_service_id: change.business_service_id || null,
  };

  showModal.value = true;
};

const saveChange = async () => {
  try {
    const payload = {
      ...form.value,
      business_service_id: form.value.business_service_id
          ? Number(form.value.business_service_id)
          : null,
    };

    if (isEditMode.value) {
      await api.put(`/changes/${editingChangeId.value}`, payload, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    } else {
      await api.post("/changes", payload, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    }

    showModal.value = false;
    resetForm();
    await fetchChanges();
  } catch (err) {
    console.error(err);
    error.value = "Impossible d'enregistrer le changement";
  }
};

const deleteChange = async (change) => {
  if (!confirm(`Supprimer ${change.reference || change.title} ?`)) return;

  try {
    await api.delete(`/changes/${change.ID}`, {
      headers: { Authorization: `Bearer ${getToken()}` },
    });

    await fetchChanges();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de supprimer le changement";
  }
};

const riskClass = (risk) => {
  return {
    Low: "badge-green",
    Medium: "badge-yellow",
    High: "badge-orange",
    Critical: "badge-red",
  }[risk] || "badge-gray";
};

const statusClass = (status) => {
  return {
    Draft: "badge-gray",
    Submitted: "badge-blue",
    Approved: "badge-green",
    Implemented: "badge-purple",
    Closed: "badge-slate",
    Rejected: "badge-red",
  }[status] || "badge-gray";
};

const submitChange = async (change) => {
  await api.put(
      `/changes/${change.ID}/submit`,
      {},
      {
        headers: {
          Authorization: `Bearer ${getToken()}`
        }
      }
  )

  await fetchChanges()
}

const approveChange = async (change) => {
  await api.put(
      `/changes/${change.ID}/approve`,
      {},
      {
        headers: {
          Authorization: `Bearer ${getToken()}`
        }
      }
  )

  await fetchChanges()
}

const rejectChange = async (change) => {
  await api.put(
      `/changes/${change.ID}/reject`,
      {},
      {
        headers: {
          Authorization: `Bearer ${getToken()}`
        }
      }
  )

  await fetchChanges()
}

const implementChange = async (change) => {
  await api.put(
      `/changes/${change.ID}/implement`,
      {},
      {
        headers: {
          Authorization: `Bearer ${getToken()}`
        }
      }
  )

  await fetchChanges()
}

const closeChange = async (change) => {
  await api.put(
      `/changes/${change.ID}/close`,
      {},
      {
        headers: {
          Authorization: `Bearer ${getToken()}`
        }
      }
  )

  await fetchChanges()
}


onMounted(async () => {
  await fetchBusinessServices();
  await fetchChanges();
});
</script>

<template>
  <DashboardLayout>
    <div class="changes-page">
      <div class="page-header">
        <div>
          <p class="eyebrow">ITSM Governance</p>
          <h1>Change Management</h1>
          <p class="subtitle">
            Gérez les changements, les risques, les plans de rollback et les services métiers impactés.
          </p>
        </div>

        <button
            v-if="role === 'Admin' || role === 'Agent'"
            class="btn-primary"
            @click="openCreateModal"
        >
          ＋ Nouveau changement
        </button>
      </div>

      <div v-if="error" class="alert-error">
        ⚠ {{ error }}
      </div>

      <div class="table-card">
        <div v-if="isLoading" class="loading">
          Chargement des changements...
        </div>

        <DataTable
            v-else
            :data="changes"
            :columns="columns"
            :options="dtOptions"
            class="display changes-table"
        >
          <template #column-reference="{ cellData }">
            <span class="reference">{{ cellData || "—" }}</span>
          </template>

          <template #column-risk="{ cellData }">
            <span class="badge" :class="riskClass(cellData)">
              {{ cellData }}
            </span>
          </template>

          <template #column-status="{ cellData }">
            <span class="badge" :class="statusClass(cellData)">
              {{ cellData }}
            </span>
          </template>

          <template #column-5="{ cellData }">
            <span>{{ cellData || "—" }}</span>
          </template>

          <template #column-6="{ rowData }">
            <div class="actions">

              <button
                  class="action-btn impact"
                  @click="router.push(`/changes/${rowData.ID}/impact`)"
                  title="Impact Analysis"
              >
                👁
              </button>

              <button
                  v-if="rowData.status === 'Draft'"
                  class="action-btn submit"
                  @click="submitChange(rowData)"
              >
                📤
              </button>

              <button
                  v-if="rowData.status === 'Submitted' && role === 'Admin'"
                  class="action-btn approve"
                  @click="approveChange(rowData)"
              >
                ✅
              </button>

              <button
                  v-if="rowData.status === 'Submitted' && role === 'Admin'"
                  class="action-btn reject"
                  @click="rejectChange(rowData)"
              >
                ❌
              </button>

              <button
                  v-if="rowData.status === 'Approved'"
                  class="action-btn implement"
                  @click="implementChange(rowData)"
              >
                🚀
              </button>

              <button
                  v-if="rowData.status === 'Implemented'"
                  class="action-btn close"
                  @click="closeChange(rowData)"
              >
                🔒
              </button>

              <button
                  class="action-btn edit"
                  @click="openEditModal(rowData)"
              >
                ✏
              </button>

              <button
                  v-if="role === 'Admin'"
                  class="action-btn delete"
                  @click="deleteChange(rowData)"
              >
                🗑
              </button>

            </div>
          </template>
        </DataTable>
      </div>
    </div>
  </DashboardLayout>

  <Teleport to="body">
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <div class="modal-header">
          <div>
            <p class="eyebrow">{{ isEditMode ? "Modification" : "Création" }}</p>
            <h2>{{ isEditMode ? "Modifier le changement" : "Nouveau changement" }}</h2>
          </div>
          <button class="close-btn" @click="showModal = false">✕</button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>Titre</label>
            <input v-model="form.title" placeholder="Migration PostgreSQL 16 vers 17" />
          </div>

          <div class="form-group">
            <label>Description</label>
            <textarea v-model="form.description" rows="3" placeholder="Décrire le changement..." />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Type</label>
              <select v-model="form.type">
                <option>Standard</option>
                <option>Normal</option>
                <option>Emergency</option>
              </select>
            </div>

            <div class="form-group">
              <label>Risque</label>
              <select v-model="form.risk">
                <option>Low</option>
                <option>Medium</option>
                <option>High</option>
                <option>Critical</option>
              </select>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>Statut</label>
              <select v-model="form.status">
                <option>Draft</option>
                <option>Submitted</option>
                <option>Approved</option>
                <option>Implemented</option>
                <option>Closed</option>
                <option>Rejected</option>
              </select>
            </div>

            <div class="form-group">
              <label>Service métier impacté</label>
              <select v-model="form.business_service_id">
                <option :value="null">Aucun</option>
                <option
                    v-for="service in businessServices"
                    :key="service.ID"
                    :value="service.ID"
                >
                  {{ service.name }}
                </option>
              </select>
            </div>
          </div>

          <div class="form-group">
            <label>Plan de rollback</label>
            <textarea
                v-model="form.rollback_plan"
                rows="3"
                placeholder="Ex: Restaurer le snapshot avant changement..."
            />
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn-ghost" @click="showModal = false">Annuler</button>
          <button class="btn-primary" @click="saveChange">
            {{ isEditMode ? "Enregistrer" : "Créer" }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.changes-page {
  padding: 2rem;
  min-height: 100vh;
  background: #0f172a;
  color: #f8fafc;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 1.5rem;
}

.eyebrow {
  color: #818cf8;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-size: 0.75rem;
  font-weight: 700;
  margin: 0 0 0.25rem;
}

h1 {
  margin: 0;
  font-size: 1.8rem;
}

.subtitle {
  color: #94a3b8;
  margin-top: 0.35rem;
  font-size: 0.9rem;
}

.table-card {
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 18px;
  overflow: hidden;
}

.loading {
  padding: 2rem;
  color: #94a3b8;
}

.alert-error {
  background: rgba(239, 68, 68, 0.12);
  border: 1px solid #ef4444;
  color: #fca5a5;
  padding: 0.75rem 1rem;
  border-radius: 10px;
  margin-bottom: 1rem;
}

.reference {
  font-family: monospace;
  color: #a5b4fc;
  font-weight: 700;
}

.badge {
  display: inline-flex;
  padding: 0.25rem 0.65rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 700;
}

.badge-green { background: rgba(34,197,94,.15); color: #86efac; }
.badge-yellow { background: rgba(234,179,8,.15); color: #fde68a; }
.badge-orange { background: rgba(249,115,22,.15); color: #fdba74; }
.badge-red { background: rgba(239,68,68,.15); color: #fca5a5; }
.badge-blue { background: rgba(59,130,246,.15); color: #93c5fd; }
.badge-purple { background: rgba(168,85,247,.15); color: #d8b4fe; }
.badge-gray,
.badge-slate { background: rgba(148,163,184,.15); color: #cbd5e1; }

.actions {
  display: flex;
  gap: 0.4rem;
}

.action-btn {
  border: none;
  background: transparent;
  color: #cbd5e1;
  cursor: pointer;
  padding: 0.35rem;
  border-radius: 8px;
}

.action-btn:hover {
  background: #334155;
}

.action-btn.delete:hover {
  background: rgba(239, 68, 68, 0.2);
}

.btn-primary {
  background: #6366f1;
  color: white;
  border: none;
  padding: 0.65rem 1.1rem;
  border-radius: 10px;
  font-weight: 700;
  cursor: pointer;
}

.btn-primary:hover {
  background: #4f46e5;
}

.btn-ghost {
  background: transparent;
  border: 1px solid #334155;
  color: #cbd5e1;
  padding: 0.65rem 1.1rem;
  border-radius: 10px;
  cursor: pointer;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(2, 6, 23, 0.75);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.modal {
  width: 100%;
  max-width: 650px;
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 18px;
  color: #f8fafc;
}

.modal-header,
.modal-footer {
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid #334155;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-footer {
  border-bottom: none;
  border-top: 1px solid #334155;
}

.close-btn {
  background: #334155;
  color: #cbd5e1;
  border: none;
  border-radius: 8px;
  padding: 0.35rem 0.6rem;
  cursor: pointer;
}

.modal-body {
  padding: 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  color: #cbd5e1;
  font-size: 0.8rem;
  font-weight: 700;
  margin-bottom: 0.4rem;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  background: #0f172a;
  color: #f8fafc;
  border: 1px solid #334155;
  border-radius: 10px;
  padding: 0.65rem 0.8rem;
}
</style>

<style>
.changes-table.dataTable {
  background: transparent !important;
  color: #e2e8f0 !important;
}

.changes-table.dataTable thead th {
  background: #1e293b !important;
  color: #94a3b8 !important;
  border-bottom: 1px solid #334155 !important;
}

.changes-table.dataTable tbody td {
  background: transparent !important;
  color: #e2e8f0 !important;
  border-bottom: 1px solid #334155 !important;
}

.changes-table.dataTable tbody tr:hover td {
  background: #243146 !important;
}

.dt-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #1e293b;
  border-bottom: 1px solid #334155;
}

.dt-footer {
  display: flex;
  justify-content: space-between;
  padding: 16px;
  background: #1e293b;
  color: #94a3b8;
}

.dt-btn {
  background: #0f172a !important;
  border: 1px solid #334155 !important;
  color: #f8fafc !important;
  border-radius: 8px !important;
  padding: 0.45rem 0.8rem !important;
  margin-right: 0.5rem !important;
}
</style>