<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute } from "vue-router";

import DataTable from "datatables.net-vue3";
import DataTablesCore from "datatables.net-dt";
import "datatables.net-dt/css/dataTables.dataTables.css";

import "datatables.net-buttons-dt";
import "datatables.net-buttons/js/buttons.html5";
import "datatables.net-buttons-dt/css/buttons.dataTables.css";

import JSZip from "jszip";

import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

window.JSZip = JSZip;
DataTable.use(DataTablesCore);

const incidents = ref([]);
const assets = ref([]);
const technicians = ref([]);

const showModal = ref(false);
const showDetailsModal = ref(false);

const selectedIncident = ref(null);
const isEditMode = ref(false);
const editingIncidentId = ref(null);

const error = ref("");
const isLoading = ref(false);
const role = localStorage.getItem("role") || "";

const highlightedIncidentId = ref(null);
const route = useRoute();

const searchQuery = ref("");
const filterStatus = ref("all");
const filterPriority = ref("all");

const comments = ref([]);
const commentMessage = ref("");
const attachments = ref([]);
const selectedFile = ref(null);

const form = ref({
  title: "",
  description: "",
  type: "",
  priority: "",
  status: "",
  asset_id: "",
  assigned_to: "",
});

const getToken = () => localStorage.getItem("token");

const priorityConfig = {
  Low: { color: "#38bdf8", bg: "rgba(56,189,248,.14)", dot: "#38bdf8", label: "Faible" },
  Medium: { color: "#a78bfa", bg: "rgba(167,139,250,.14)", dot: "#a78bfa", label: "Moyenne" },
  High: { color: "#f59e0b", bg: "rgba(245,158,11,.16)", dot: "#f59e0b", label: "Haute" },
  Critical: { color: "#fb7185", bg: "rgba(251,113,133,.16)", dot: "#fb7185", label: "Critique" },
};

const statusConfig = {
  Open: { color: "#38bdf8", bg: "rgba(56,189,248,.14)", dot: "#38bdf8", label: "Ouvert" },
  "In Progress": { color: "#a78bfa", bg: "rgba(167,139,250,.14)", dot: "#a78bfa", label: "En cours" },
  Resolved: { color: "#34d399", bg: "rgba(52,211,153,.14)", dot: "#34d399", label: "Résolu" },
  Closed: { color: "#94a3b8", bg: "rgba(148,163,184,.14)", dot: "#94a3b8", label: "Clôturé" },
};

const typeOptions = ["Panne", "Réseau", "Sécurité", "Matériel", "Logiciel", "Accès", "Autre"];

const headers = () => ({
  headers: { Authorization: `Bearer ${getToken()}` },
});

const getBadgeStyle = (map, key) => {
  const cfg = map[key];
  if (!cfg) return {};
  return {
    color: cfg.color,
    backgroundColor: cfg.bg,
    "--dot": cfg.dot,
  };
};

const getSlaStatus = (incident) => {
  if (!incident) {
    return { label: "Non défini", color: "#94a3b8", bg: "rgba(148,163,184,.14)", icon: "—" };
  }

  if (incident.status === "Closed" || incident.status === "Resolved") {
    return { label: "Terminé", color: "#94a3b8", bg: "rgba(148,163,184,.14)", icon: "✓" };
  }

  if (!incident.due_at) {
    return { label: "Non défini", color: "#94a3b8", bg: "rgba(148,163,184,.14)", icon: "—" };
  }

  const diffH = (new Date(incident.due_at) - new Date()) / 3600000;

  if (diffH < 0) {
    return { label: "En retard", color: "#fb7185", bg: "rgba(251,113,133,.16)", icon: "⚠" };
  }

  if (diffH <= 24) {
    return { label: "Échéance proche", color: "#f59e0b", bg: "rgba(245,158,11,.16)", icon: "⏱" };
  }

  return { label: "Dans les délais", color: "#34d399", bg: "rgba(52,211,153,.14)", icon: "✓" };
};

const formatDate = (date) => {
  if (!date) return "—";

  return new Date(date).toLocaleDateString("fr-FR", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const initials = (name) =>
    name
        ? name
            .split(" ")
            .map((part) => part[0])
            .join("")
            .toUpperCase()
            .slice(0, 2)
        : "?";

const stats = computed(() => ({
  total: incidents.value.length,
  open: incidents.value.filter((item) => item.status === "Open").length,
  inProgress: incidents.value.filter((item) => item.status === "In Progress").length,
  critical: incidents.value.filter(
      (item) => item.priority === "Critical" && item.status !== "Closed"
  ).length,
  overdue: incidents.value.filter((item) => {
    if (!item.due_at || item.status === "Closed" || item.status === "Resolved") {
      return false;
    }

    return new Date(item.due_at) < new Date();
  }).length,
}));

const filteredIncidents = computed(() => {
  let list = incidents.value;

  if (filterStatus.value !== "all") {
    list = list.filter((item) => item.status === filterStatus.value);
  }

  if (filterPriority.value !== "all") {
    list = list.filter((item) => item.priority === filterPriority.value);
  }

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();

    list = list.filter(
        (item) =>
            item.title?.toLowerCase().includes(query) ||
            item.type?.toLowerCase().includes(query) ||
            item.asset?.name?.toLowerCase().includes(query) ||
            item.assigned_to?.toLowerCase().includes(query) ||
            String(item.ID || "").includes(query)
    );
  }

  return list;
});

const datatableRows = computed(() =>
    filteredIncidents.value.map((item) => ({
      id: item.ID,
      incident: item.title,
      type: item.type || "—",
      priority: item.priority,
      status: item.status,
      sla: getSlaStatus(item).label,
      asset: item.asset?.name || "Non lié",
      assigned_to: item.assigned_to || "Non assigné",
      created_at: formatDate(item.CreatedAt || item.created_at),
      raw: item,
    }))
);

const columns = [
  { data: "incident", title: "Incident", render: "#incident" },
  { data: "priority", title: "Priorité", render: "#priority" },
  { data: "status", title: "Statut", render: "#status" },
  { data: "sla", title: "SLA", render: "#sla" },
  { data: "asset", title: "Asset lié", render: "#asset" },
  { data: "assigned_to", title: "Assigné à", render: "#assigned" },
  { data: "created_at", title: "Créé le" },
  {
    data: null,
    title: "",
    orderable: false,
    searchable: false,
    render: "#actions",
  },
];

const tableOptions = {
  dom: '<"dt-top"Bf>rt<"dt-bottom"lip>',
  pageLength: 25,
  lengthMenu: [10, 25, 50, 100],
  order: [[0, "desc"]],
  autoWidth: false,
  stateSave: true,
  buttons: [
    {
      extend: "csvHtml5",
      text: "Exporter CSV",
      className: "dt-export-btn",
      filename: "secureasset_incidents",
      exportOptions: {
        columns: [0, 1, 2, 3, 4, 5, 6],
      },
    },
    {
      extend: "excelHtml5",
      text: "Exporter Excel",
      className: "dt-export-btn",
      filename: "secureasset_incidents",
      title: "SecureAsset Manager - Incidents",
      exportOptions: {
        columns: [0, 1, 2, 3, 4, 5, 6],
      },
    },
  ],
  language: {
    search: "Recherche DataTable :",
    lengthMenu: "Afficher _MENU_ lignes",
    info: "_START_ à _END_ sur _TOTAL_ incidents",
    infoEmpty: "Aucun incident",
    infoFiltered: "(filtré sur _MAX_ incidents)",
    zeroRecords: "Aucun incident trouvé",
    paginate: {
      first: "Premier",
      last: "Dernier",
      next: "Suivant",
      previous: "Précédent",
    },
  },
};

const resetForm = () => {
  form.value = {
    title: "",
    description: "",
    type: "",
    priority: "",
    status: "",
    asset_id: "",
    assigned_to: "",
  };

  isEditMode.value = false;
  editingIncidentId.value = null;
  error.value = "";
};

const fetchTechnicians = async () => {
  try {
    const response = await api.get("/technicians", headers());
    technicians.value = response.data || [];
  } catch (err) {
    console.error(err);
  }
};

const fetchIncidents = async () => {
  isLoading.value = true;

  try {
    const response = await api.get("/incidents", headers());
    incidents.value = response.data || [];
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger les incidents";
  } finally {
    isLoading.value = false;
  }
};

const fetchAssets = async () => {
  try {
    const response = await api.get("/assets", headers());
    assets.value = response.data || [];
  } catch (err) {
    console.error(err);
  }
};

const openCreateModal = () => {
  resetForm();
  form.value.status = "Open";
  showModal.value = true;
};

const openEditModal = (incident) => {
  isEditMode.value = true;
  editingIncidentId.value = incident.ID;

  form.value = {
    title: incident.title || "",
    description: incident.description || "",
    type: incident.type || "",
    priority: incident.priority || "",
    status: incident.status || "Open",
    asset_id: incident.asset_id || "",
    assigned_to: incident.assigned_to || "",
  };

  showModal.value = true;
};

const openDetailsModal = async (incident) => {
  selectedIncident.value = incident;
  comments.value = [];
  attachments.value = [];
  showDetailsModal.value = true;

  await fetchComments(incident.ID);
  await fetchAttachments(incident.ID);
};

const saveIncident = async () => {
  try {
    error.value = "";

    const payload = {
      ...form.value,
      asset_id: form.value.asset_id ? Number(form.value.asset_id) : null,
    };

    if (isEditMode.value) {
      await api.put(`/incidents/${editingIncidentId.value}`, payload, headers());
    } else {
      await api.post("/incidents", payload, headers());
    }

    showModal.value = false;
    resetForm();
    await fetchIncidents();
  } catch (err) {
    console.error(err);
    error.value = "Erreur lors de l'enregistrement de l'incident";
  }
};

const deleteIncident = async (incident) => {
  if (!confirm(`Supprimer l'incident « ${incident.title} » ?`)) return;

  try {
    await api.delete(`/incidents/${incident.ID}`, headers());
    await fetchIncidents();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de supprimer cet incident";
  }
};

const takeIncident = async (incident) => {
  try {
    await api.put(`/incidents/${incident.ID}/take`, {}, headers());
    await fetchIncidents();
  } catch (err) {
    console.error(err);
    error.value = "Impossible d'assigner le ticket";
  }
};

const resolveIncident = async (incident) => {
  try {
    await api.put(`/incidents/${incident.ID}/resolve`, {}, headers());
    await fetchIncidents();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de résoudre le ticket";
  }
};

const closeIncident = async (incident) => {
  try {
    await api.put(`/incidents/${incident.ID}/close`, {}, headers());
    await fetchIncidents();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de clôturer le ticket";
  }
};

const fetchComments = async (id) => {
  try {
    const response = await api.get(`/incidents/${id}/comments`, headers());
    comments.value = response.data || [];
  } catch (err) {
    console.error(err);
    comments.value = [];
  }
};

const addComment = async () => {
  if (!commentMessage.value.trim()) return;

  try {
    await api.post(
        `/incidents/${selectedIncident.value.ID}/comments`,
        { message: commentMessage.value },
        headers()
    );

    commentMessage.value = "";
    await fetchComments(selectedIncident.value.ID);
  } catch (err) {
    console.error(err);
    error.value = "Impossible d'ajouter le commentaire";
  }
};

const fetchAttachments = async (id) => {
  try {
    const response = await api.get(`/incidents/${id}/attachments`, headers());
    attachments.value = response.data || [];
  } catch (err) {
    console.error(err);
    attachments.value = [];
  }
};

const onFileChange = (event) => {
  selectedFile.value = event.target.files[0];
};

const uploadAttachment = async () => {
  if (!selectedFile.value) return;

  try {
    const formData = new FormData();
    formData.append("file", selectedFile.value);

    await api.post(`/incidents/${selectedIncident.value.ID}/attachments`, formData, {
      headers: {
        Authorization: `Bearer ${getToken()}`,
        "Content-Type": "multipart/form-data",
      },
    });

    selectedFile.value = null;
    await fetchAttachments(selectedIncident.value.ID);
  } catch (err) {
    console.error(err);
    error.value = "Impossible d'uploader la pièce jointe";
  }
};

const downloadAttachment = async (attachment) => {
  try {
    const response = await api.get(`/attachments/${attachment.ID}/download`, {
      ...headers(),
      responseType: "blob",
    });

    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement("a");

    link.href = url;
    link.setAttribute("download", attachment.file_name);

    document.body.appendChild(link);
    link.click();
    link.remove();

    window.URL.revokeObjectURL(url);
  } catch (err) {
    console.error(err);
    error.value = "Impossible de télécharger la pièce jointe";
  }
};

const deleteAttachment = async (attachment) => {
  if (!confirm(`Supprimer ${attachment.file_name} ?`)) return;

  try {
    await api.delete(`/attachments/${attachment.ID}`, headers());
    await fetchAttachments(selectedIncident.value.ID);
  } catch (err) {
    console.error(err);
    error.value = "Impossible de supprimer la pièce jointe";
  }
};

onMounted(async () => {
  await fetchIncidents();
  await fetchTechnicians();
  await fetchAssets();

  highlightedIncidentId.value = Number(route.query.incident);

  const id = route.query.incident;

  if (id) {
    const incident = incidents.value.find((item) => item.ID === Number(id));

    if (incident) {
      await openDetailsModal(incident);
    }
  }
});
</script>

<template>
  <DashboardLayout>
    <div class="incidents-page">
      <div class="page-header">
        <div>
          <p class="page-eyebrow">Gestion ITSM</p>
          <h1 class="page-title">Incidents</h1>
        </div>

        <button
            v-if="role === 'Requester'"
            class="btn-primary"
            @click="openCreateModal"
        >
          <span>＋</span>
          Nouvel incident
        </button>
      </div>

      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon neutral">🎫</div>
          <div>
            <div class="stat-value">{{ stats.total }}</div>
            <div class="stat-label">Total</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon blue">📂</div>
          <div>
            <div class="stat-value blue">{{ stats.open }}</div>
            <div class="stat-label">Ouverts</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon purple">⚙</div>
          <div>
            <div class="stat-value purple">{{ stats.inProgress }}</div>
            <div class="stat-label">En cours</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon red">🔴</div>
          <div>
            <div class="stat-value red">{{ stats.critical }}</div>
            <div class="stat-label">Critiques</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon amber">⏰</div>
          <div>
            <div class="stat-value amber">{{ stats.overdue }}</div>
            <div class="stat-label">En retard</div>
          </div>
        </div>
      </div>

      <div class="toolbar">
        <div class="search-wrap">
          <span class="search-icon">🔍</span>
          <input
              v-model="searchQuery"
              type="text"
              placeholder="Rechercher titre, type, asset, ID…"
              class="search-input"
          />
        </div>

        <div class="filter-group">
          <select v-model="filterStatus" class="filter-select">
            <option value="all">Tous les statuts</option>
            <option
                v-for="(value, key) in statusConfig"
                :key="key"
                :value="key"
            >
              {{ value.label }}
            </option>
          </select>

          <select v-model="filterPriority" class="filter-select">
            <option value="all">Toutes priorités</option>
            <option
                v-for="(value, key) in priorityConfig"
                :key="key"
                :value="key"
            >
              {{ value.label }}
            </option>
          </select>
        </div>
      </div>

      <div v-if="error" class="alert-error">
        <span>⚠</span>
        {{ error }}
      </div>

      <div class="table-card">
        <div v-if="isLoading" class="table-loading">
          <span class="spinner"></span>
          Chargement des incidents…
        </div>

        <DataTable
            v-else
            :data="datatableRows"
            :columns="columns"
            :options="tableOptions"
            class="display inc-datatable"
        >
          <template #incident="{ rowData }">
            <div class="cell-incident">
              <div class="incident-title">{{ rowData.raw.title }}</div>
              <div class="incident-meta">
                <span class="type-tag">{{ rowData.raw.type || "—" }}</span>
                <span class="incident-id">#{{ rowData.raw.ID }}</span>
              </div>
            </div>
          </template>

          <template #priority="{ rowData }">
            <span
                class="badge"
                :style="getBadgeStyle(priorityConfig, rowData.raw.priority)"
            >
              <span class="badge-dot"></span>
              {{ priorityConfig[rowData.raw.priority]?.label || rowData.raw.priority }}
            </span>
          </template>

          <template #status="{ rowData }">
            <span
                class="badge"
                :style="getBadgeStyle(statusConfig, rowData.raw.status)"
            >
              <span class="badge-dot"></span>
              {{ statusConfig[rowData.raw.status]?.label || rowData.raw.status }}
            </span>
          </template>

          <template #sla="{ rowData }">
            <div
                class="sla-cell"
                :style="{ color: getSlaStatus(rowData.raw).color }"
            >
              <span class="sla-icon">{{ getSlaStatus(rowData.raw).icon }}</span>
              <span class="sla-label">{{ getSlaStatus(rowData.raw).label }}</span>
            </div>

            <div v-if="rowData.raw.due_at" class="sla-date">
              {{ formatDate(rowData.raw.due_at) }}
            </div>
          </template>

          <template #asset="{ rowData }">
            <span v-if="rowData.raw.asset?.name" class="asset-chip">
              🖥 {{ rowData.raw.asset.name }}
            </span>
            <span v-else class="cell-muted">Non lié</span>
          </template>

          <template #assigned="{ rowData }">
            <div v-if="rowData.raw.assigned_to" class="assigned-user">
              <span class="user-avatar">{{ initials(rowData.raw.assigned_to) }}</span>
              <span class="user-name">{{ rowData.raw.assigned_to }}</span>
            </div>

            <span v-else class="cell-muted unassigned">Non assigné</span>
          </template>

          <template #actions="{ rowData }">
            <div class="cell-actions">
              <button
                  class="action-btn"
                  title="Voir les détails"
                  @click="openDetailsModal(rowData.raw)"
              >
                🔍
              </button>

              <button
                  v-if="(role === 'Admin' || role === 'Requester') && rowData.raw.status !== 'Closed'"
                  class="action-btn edit"
                  title="Modifier"
                  @click="openEditModal(rowData.raw)"
              >
                ✏
              </button>

              <button
                  v-if="(role === 'Admin' || role === 'Requester') && rowData.raw.status !== 'Closed'"
                  class="action-btn del"
                  title="Supprimer"
                  @click="deleteIncident(rowData.raw)"
              >
                🗑
              </button>

              <button
                  v-if="role === 'Agent' && rowData.raw.status !== 'Resolved' && rowData.raw.status !== 'Closed'"
                  class="action-btn take"
                  title="Prendre en charge"
                  @click="takeIncident(rowData.raw)"
              >
                ⚡
              </button>

              <button
                  v-if="role === 'Agent' && rowData.raw.status === 'In Progress'"
                  class="action-btn resolve"
                  title="Résoudre"
                  @click="resolveIncident(rowData.raw)"
              >
                ✅
              </button>

              <button
                  v-if="role === 'Agent' && rowData.raw.status === 'Resolved'"
                  class="action-btn close-btn"
                  title="Clôturer"
                  @click="closeIncident(rowData.raw)"
              >
                🔒
              </button>
            </div>
          </template>
        </DataTable>
      </div>
    </div>
  </DashboardLayout>

  <Teleport to="body">
    <div
        v-if="showModal"
        class="modal-overlay"
        @click.self="showModal = false"
    >
      <div class="modal">
        <div class="modal-header">
          <div>
            <p class="modal-eyebrow">
              {{ isEditMode ? "Modification" : "Création" }}
            </p>

            <h2 class="modal-title">
              {{ isEditMode ? "Modifier l'incident" : "Nouvel incident" }}
            </h2>
          </div>

          <button class="modal-close" @click="showModal = false">✕</button>
        </div>

        <div v-if="error" class="alert-error modal-alert">
          <span>⚠</span>
          {{ error }}
        </div>

        <div class="modal-body">
          <div class="form-section">
            <p class="form-section-label">Identification</p>

            <div class="form-group full">
              <label>Titre <span class="required">*</span></label>
              <input
                  v-model="form.title"
                  type="text"
                  placeholder="Résumé concis de l'incident"
              />
            </div>

            <div class="form-group full">
              <label>Description</label>
              <textarea
                  v-model="form.description"
                  rows="3"
                  placeholder="Contexte, symptômes, actions déjà tentées…"
              ></textarea>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label>Type</label>
                <select v-model="form.type">
                  <option value="" disabled>Sélectionner…</option>
                  <option
                      v-for="type in typeOptions"
                      :key="type"
                      :value="type"
                  >
                    {{ type }}
                  </option>
                </select>
              </div>

              <div class="form-group">
                <label>Priorité <span class="required">*</span></label>
                <select v-model="form.priority">
                  <option value="" disabled>Sélectionner…</option>
                  <option value="Low">Faible</option>
                  <option value="Medium">Moyenne</option>
                  <option value="High">Haute</option>
                  <option value="Critical">Critique</option>
                </select>
              </div>
            </div>
          </div>

          <div class="form-section">
            <p class="form-section-label">Affectation</p>

            <div class="form-row">
              <div class="form-group">
                <label>Statut</label>
                <select v-if="role !== 'Requester'" v-model="form.status">
                  <option value="" disabled>Sélectionner…</option>
                  <option value="Open">Ouvert</option>
                  <option value="In Progress">En cours</option>
                  <option value="Resolved">Résolu</option>
                  <option value="Closed">Clôturé</option>
                </select>

                <input v-else :value="'Ouvert'" disabled />
              </div>

              <div class="form-group">
                <label>Technicien</label>
                <select v-model="form.assigned_to">
                  <option value="">Non assigné</option>
                  <option
                      v-for="technician in technicians"
                      :key="technician.ID"
                      :value="technician.email"
                  >
                    {{ technician.name }} — {{ technician.email }}
                  </option>
                </select>
              </div>
            </div>

            <div class="form-group full">
              <label>Asset concerné</label>
              <select v-model="form.asset_id">
                <option value="">Aucun asset lié</option>
                <option
                    v-for="asset in assets"
                    :key="asset.ID"
                    :value="asset.ID"
                >
                  {{ asset.name }} — {{ asset.ip_address }}
                </option>
              </select>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn-ghost" @click="showModal = false">
            Annuler
          </button>

          <button class="btn-primary" @click="saveIncident">
            {{ isEditMode ? "Enregistrer les modifications" : "Créer l'incident" }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>

  <Teleport to="body">
    <div
        v-if="showDetailsModal"
        class="modal-overlay"
        @click.self="showDetailsModal = false"
    >
      <div class="modal modal-lg">
        <div class="modal-header">
          <div style="flex: 1; min-width: 0">
            <p class="modal-eyebrow">Incident #{{ selectedIncident?.ID }}</p>
            <h2 class="modal-title truncate">{{ selectedIncident?.title }}</h2>
          </div>

          <button class="modal-close" @click="showDetailsModal = false">✕</button>
        </div>

        <div class="details-body">
          <div class="details-main">
            <div class="details-badges">
              <span
                  class="badge"
                  :style="getBadgeStyle(priorityConfig, selectedIncident?.priority)"
              >
                <span class="badge-dot"></span>
                {{ priorityConfig[selectedIncident?.priority]?.label || selectedIncident?.priority }}
              </span>

              <span
                  class="badge"
                  :style="getBadgeStyle(statusConfig, selectedIncident?.status)"
              >
                <span class="badge-dot"></span>
                {{ statusConfig[selectedIncident?.status]?.label || selectedIncident?.status }}
              </span>

              <span
                  class="badge sla-badge"
                  :style="{
                  color: getSlaStatus(selectedIncident).color,
                  backgroundColor: getSlaStatus(selectedIncident).bg,
                }"
              >
                {{ getSlaStatus(selectedIncident).icon }}
                {{ getSlaStatus(selectedIncident).label }}
              </span>
            </div>

            <div class="detail-block">
              <p class="detail-block-label">Description</p>
              <p class="detail-block-text">
                {{ selectedIncident?.description || "Aucune description fournie." }}
              </p>
            </div>

            <div class="detail-meta-grid">
              <div class="meta-item">
                <span class="meta-label">Type</span>
                <span class="meta-value">{{ selectedIncident?.type || "—" }}</span>
              </div>

              <div class="meta-item">
                <span class="meta-label">Asset</span>
                <span class="meta-value">{{ selectedIncident?.asset?.name || "Non lié" }}</span>
              </div>

              <div class="meta-item">
                <span class="meta-label">Assigné à</span>
                <span class="meta-value">{{ selectedIncident?.assigned_to || "Non assigné" }}</span>
              </div>

              <div class="meta-item">
                <span class="meta-label">Échéance SLA</span>
                <span class="meta-value">{{ formatDate(selectedIncident?.due_at) }}</span>
              </div>

              <div class="meta-item">
                <span class="meta-label">Créé le</span>
                <span class="meta-value">{{ formatDate(selectedIncident?.CreatedAt) }}</span>
              </div>

              <div class="meta-item">
                <span class="meta-label">Dernière MàJ</span>
                <span class="meta-value">{{ formatDate(selectedIncident?.UpdatedAt) }}</span>
              </div>
            </div>

            <div class="detail-block">
              <p class="detail-block-label">
                Historique & commentaires ({{ comments.length }})
              </p>

              <div class="comments-list">
                <div
                    v-for="comment in comments"
                    :key="comment.ID"
                    class="comment-item"
                >
                  <div class="comment-avatar">
                    {{ initials(comment.user?.name || "Sys") }}
                  </div>

                  <div class="comment-body">
                    <div class="comment-meta">
                      <span class="comment-author">
                        {{ comment.user?.name || "Système" }}
                      </span>

                      <span class="comment-type-badge">{{ comment.type }}</span>

                      <span class="comment-date">
                        {{ formatDate(comment.CreatedAt) }}
                      </span>
                    </div>

                    <p class="comment-text">{{ comment.message }}</p>
                  </div>
                </div>

                <p v-if="comments.length === 0" class="empty-comments">
                  Aucun commentaire pour ce ticket.
                </p>
              </div>

              <div
                  v-if="selectedIncident?.status !== 'Closed' && role !== 'Admin'"
                  class="comment-form"
              >
                <textarea
                    v-model="commentMessage"
                    rows="2"
                    placeholder="Ajouter un commentaire…"
                ></textarea>

                <button class="btn-primary btn-sm" @click="addComment">
                  Publier
                </button>
              </div>
            </div>
          </div>

          <div class="details-sidebar">
            <div class="detail-block">
              <p class="detail-block-label">
                Pièces jointes ({{ attachments.length }})
              </p>

              <div class="attachment-list">
                <div
                    v-for="attachment in attachments"
                    :key="attachment.ID"
                    class="attachment-item"
                >
                  <div class="attachment-info">
                    <span class="attachment-icon">📎</span>

                    <div>
                      <div class="attachment-name">
                        {{ attachment.file_name }}
                      </div>

                      <div class="attachment-size">
                        {{ Math.round(attachment.file_size / 1024) }} KB
                      </div>
                    </div>
                  </div>

                  <div class="attachment-actions">
                    <button
                        class="link-btn blue"
                        @click="downloadAttachment(attachment)"
                    >
                      ↓
                    </button>

                    <button
                        v-if="selectedIncident?.status !== 'Closed'"
                        class="link-btn red"
                        @click="deleteAttachment(attachment)"
                    >
                      ✕
                    </button>
                  </div>
                </div>

                <p v-if="attachments.length === 0" class="empty-comments">
                  Aucune pièce jointe.
                </p>
              </div>

              <div
                  v-if="selectedIncident?.status !== 'Closed'"
                  class="upload-zone"
              >
                <input
                    id="file-upload"
                    type="file"
                    class="file-input"
                    @change="onFileChange"
                />

                <label for="file-upload" class="file-label">
                  <span>📁</span>
                  <span>
                    {{ selectedFile ? selectedFile.name : "Choisir un fichier…" }}
                  </span>
                </label>

                <button
                    v-if="selectedFile"
                    class="btn-primary btn-sm"
                    @click="uploadAttachment"
                >
                  Uploader
                </button>
              </div>
            </div>

            <div v-if="role === 'Agent'" class="detail-block">
              <p class="detail-block-label">Actions rapides</p>

              <div class="quick-actions">
                <button
                    v-if="selectedIncident?.status !== 'Resolved' && selectedIncident?.status !== 'Closed'"
                    class="quick-btn indigo"
                    @click="takeIncident(selectedIncident); showDetailsModal = false"
                >
                  ⚡ Prendre en charge
                </button>

                <button
                    v-if="selectedIncident?.status === 'In Progress'"
                    class="quick-btn green"
                    @click="resolveIncident(selectedIncident); showDetailsModal = false"
                >
                  ✅ Marquer résolu
                </button>

                <button
                    v-if="selectedIncident?.status === 'Resolved'"
                    class="quick-btn gray"
                    @click="closeIncident(selectedIncident); showDetailsModal = false"
                >
                  🔒 Clôturer
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn-ghost" @click="showDetailsModal = false">
            Fermer
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.incidents-page {
  min-height: 100vh;
  padding: 2rem;
  max-width: 1480px;
  color: #e5e7eb;
  font-family: Inter, "Segoe UI", system-ui, sans-serif;
}

.page-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  margin-bottom: 1.75rem;
}

.page-eyebrow {
  margin: 0 0 0.25rem;
  color: #a78bfa;
  font-size: 0.72rem;
  font-weight: 700;
  letter-spacing: 0.09em;
  text-transform: uppercase;
}

.page-title {
  margin: 0;
  color: #f8fafc;
  font-size: 1.75rem;
  font-weight: 800;
  letter-spacing: -0.02em;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 0.875rem;
  padding: 1rem 1.125rem;
  background: #111827;
  border: 1px solid #1e293b;
  border-radius: 14px;
  box-shadow: 0 14px 30px rgba(0, 0, 0, 0.22);
}

.stat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border-radius: 10px;
  font-size: 1.125rem;
  flex-shrink: 0;
}

.stat-icon.neutral { background: rgba(148, 163, 184, 0.12); }
.stat-icon.blue { background: rgba(56, 189, 248, 0.14); }
.stat-icon.purple { background: rgba(167, 139, 250, 0.14); }
.stat-icon.red { background: rgba(251, 113, 133, 0.14); }
.stat-icon.amber { background: rgba(245, 158, 11, 0.14); }

.stat-value {
  color: #f8fafc;
  font-size: 1.5rem;
  font-weight: 800;
  line-height: 1;
}

.stat-value.blue { color: #38bdf8; }
.stat-value.purple { color: #a78bfa; }
.stat-value.red { color: #fb7185; }
.stat-value.amber { color: #f59e0b; }

.stat-label {
  margin-top: 0.1rem;
  color: #94a3b8;
  font-size: 0.78rem;
  font-weight: 600;
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.25rem;
}

.search-wrap {
  position: relative;
  min-width: 260px;
  flex: 1;
}

.search-icon {
  position: absolute;
  top: 50%;
  left: 0.875rem;
  transform: translateY(-50%);
  font-size: 0.9rem;
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 0.7rem 0.875rem 0.7rem 2.5rem;
  background: #111827;
  color: #f8fafc;
  border: 1px solid #334155;
  border-radius: 12px;
  outline: none;
}

.search-input::placeholder {
  color: #64748b;
}

.search-input:focus {
  border-color: #7c3aed;
  box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.18);
}

.filter-group {
  display: flex;
  gap: 0.5rem;
}

.filter-select {
  min-width: 170px;
  padding: 0.7rem 2rem 0.7rem 0.875rem;
  background: #111827;
  color: #f8fafc;
  border: 1px solid #334155;
  border-radius: 12px;
  outline: none;
  cursor: pointer;
}

.filter-select option {
  background: #0f172a;
  color: #f8fafc;
}

.alert-error {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
  padding: 0.75rem 1rem;
  background: rgba(239, 68, 68, 0.14);
  border: 1px solid rgba(239, 68, 68, 0.32);
  color: #fca5a5;
  border-radius: 12px;
  font-size: 0.875rem;
}

.table-card {
  padding: 1rem;
  background: #111827;
  border: 1px solid #1e293b;
  border-radius: 18px;
  overflow: hidden;
  box-shadow: 0 16px 36px rgba(0, 0, 0, 0.24);
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
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 2px solid #334155;
  border-top-color: #a78bfa;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

:deep(.inc-datatable) {
  width: 100% !important;
  color: #e5e7eb;
  border-collapse: collapse !important;
}

:deep(.dt-container),
:deep(.dataTables_wrapper) {
  color: #e5e7eb;
}

:deep(.dt-top) {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

:deep(.dt-buttons) {
  display: flex;
  gap: 0.5rem;
}

:deep(.dt-button),
:deep(.dt-export-btn) {
  padding: 0.55rem 0.9rem !important;
  background: #1e293b !important;
  color: #f8fafc !important;
  border: 1px solid #334155 !important;
  border-radius: 10px !important;
  font-size: 0.8rem !important;
  font-weight: 700 !important;
  cursor: pointer !important;
  box-shadow: none !important;
}

:deep(.dt-button:hover),
:deep(.dt-export-btn:hover) {
  background: #334155 !important;
  border-color: #475569 !important;
}

:deep(.dt-search) {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #94a3b8;
  font-size: 0.82rem;
}

:deep(.dt-search input) {
  min-width: 240px;
  background: #0f172a !important;
  color: #f8fafc !important;
  border: 1px solid #334155 !important;
  border-radius: 10px !important;
  padding: 0.55rem 0.75rem !important;
  outline: none !important;
}

:deep(.dt-length) {
  color: #94a3b8;
  font-size: 0.82rem;
}

:deep(.dt-length select) {
  background: #0f172a !important;
  color: #f8fafc !important;
  border: 1px solid #334155 !important;
  border-radius: 8px !important;
  padding: 0.35rem 0.55rem !important;
}

:deep(table.dataTable thead th) {
  padding: 0.85rem 1rem !important;
  background: #0f172a !important;
  color: #94a3b8 !important;
  border-bottom: 1px solid #334155 !important;
  font-size: 0.72rem !important;
  font-weight: 800 !important;
  letter-spacing: 0.06em !important;
  text-transform: uppercase !important;
  white-space: nowrap !important;
}

:deep(table.dataTable tbody td) {
  padding: 0.9rem 1rem !important;
  background: #111827 !important;
  color: #cbd5e1 !important;
  border-top: 1px solid #1e293b !important;
  vertical-align: middle !important;
}

:deep(table.dataTable tbody tr:hover td) {
  background: #172033 !important;
}

:deep(.dt-bottom) {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  margin-top: 1rem;
  color: #94a3b8;
  font-size: 0.82rem;
}

:deep(.dt-paging) {
  display: flex;
  gap: 0.35rem;
}

:deep(.dt-paging-button) {
  min-width: 34px;
  padding: 0.45rem 0.7rem !important;
  background: #0f172a !important;
  color: #cbd5e1 !important;
  border: 1px solid #334155 !important;
  border-radius: 8px !important;
  cursor: pointer !important;
}

:deep(.dt-paging-button.current) {
  background: #7c3aed !important;
  color: white !important;
  border-color: #7c3aed !important;
}

:deep(.dt-paging-button.disabled) {
  opacity: 0.4 !important;
  cursor: not-allowed !important;
}

.cell-incident .incident-title {
  margin-bottom: 0.25rem;
  color: #f8fafc;
  font-weight: 800;
}

.incident-meta {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.type-tag {
  padding: 0.15rem 0.5rem;
  background: rgba(148, 163, 184, 0.14);
  color: #cbd5e1;
  border-radius: 6px;
  font-size: 0.7rem;
  font-weight: 700;
}

.incident-id {
  color: #64748b;
  font-size: 0.72rem;
}

.badge {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.25rem 0.65rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 800;
  white-space: nowrap;
}

.badge-dot {
  width: 6px;
  height: 6px;
  background: var(--dot, currentColor);
  border-radius: 50%;
  flex-shrink: 0;
}

.sla-cell {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.8rem;
  font-weight: 800;
}

.sla-date {
  margin-top: 0.15rem;
  color: #64748b;
  font-size: 0.7rem;
}

.asset-chip {
  color: #cbd5e1;
  font-size: 0.82rem;
}

.cell-muted {
  color: #64748b;
  font-size: 0.85rem;
}

.unassigned {
  font-style: italic;
}

.assigned-user {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.user-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  background: rgba(167, 139, 250, 0.16);
  color: #c4b5fd;
  border-radius: 999px;
  font-size: 0.7rem;
  font-weight: 900;
  flex-shrink: 0;
}

.user-name {
  color: #cbd5e1;
  font-size: 0.82rem;
}

.cell-actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 0.25rem;
}

.action-btn {
  padding: 0.4rem;
  background: transparent;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  opacity: 0.75;
  transition: opacity 0.15s, background 0.15s;
}

.action-btn:hover {
  opacity: 1;
}

.action-btn.edit:hover { background: rgba(59, 130, 246, 0.18); }
.action-btn.del:hover { background: rgba(239, 68, 68, 0.18); }
.action-btn.take:hover { background: rgba(124, 58, 237, 0.18); }
.action-btn.resolve:hover { background: rgba(34, 197, 94, 0.18); }
.action-btn.close-btn:hover { background: rgba(148, 163, 184, 0.18); }

.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: #7c3aed;
  color: white;
  border: none;
  border-radius: 12px;
  padding: 0.68rem 1.1rem;
  font-size: 0.875rem;
  font-weight: 800;
  cursor: pointer;
  transition: background 0.15s, box-shadow 0.15s, transform 0.1s;
}

.btn-primary:hover {
  background: #6d28d9;
  box-shadow: 0 8px 22px rgba(124, 58, 237, 0.32);
}

.btn-primary:active {
  transform: scale(0.98);
}

.btn-sm {
  padding: 0.48rem 0.875rem;
  font-size: 0.8rem;
}

.btn-ghost {
  padding: 0.65rem 1.1rem;
  background: transparent;
  color: #cbd5e1;
  border: 1px solid #334155;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 700;
  cursor: pointer;
}

.btn-ghost:hover {
  background: #1e293b;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(2, 6, 23, 0.75);
  backdrop-filter: blur(4px);
}

.modal {
  width: 100%;
  max-width: 580px;
  max-height: 90vh;
  overflow-y: auto;
  background: #111827;
  color: #e5e7eb;
  border: 1px solid #1e293b;
  border-radius: 18px;
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.42);
}

.modal-lg {
  max-width: 920px;
}

.modal-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 1.5rem 1.75rem 1.25rem;
  border-bottom: 1px solid #1e293b;
}

.modal-eyebrow {
  margin: 0 0 0.2rem;
  color: #a78bfa;
  font-size: 0.7rem;
  font-weight: 800;
  letter-spacing: 0.1em;
  text-transform: uppercase;
}

.modal-title {
  margin: 0;
  color: #f8fafc;
  font-size: 1.2rem;
  font-weight: 800;
}

.truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.modal-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  background: #1e293b;
  color: #94a3b8;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

.modal-close:hover {
  background: #334155;
  color: #f8fafc;
}

.modal-body {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  padding: 1.25rem 1.75rem;
}

.modal-alert {
  margin-bottom: 0;
  border-radius: 0;
  border-left: none;
  border-right: none;
  border-top: none;
}

.form-section {
  display: flex;
  flex-direction: column;
}

.form-section-label {
  margin: 0 0 0.875rem;
  color: #94a3b8;
  font-size: 0.72rem;
  font-weight: 900;
  letter-spacing: 0.08em;
  text-transform: uppercase;
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

.form-group.full {
  grid-column: 1 / -1;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-group label {
  color: #cbd5e1;
  font-size: 0.8rem;
  font-weight: 800;
}

.required {
  color: #fb7185;
}

.form-group input,
.form-group select,
.form-group textarea,
.comment-form textarea {
  padding: 0.65rem 0.875rem;
  background: #0f172a;
  color: #f8fafc;
  border: 1px solid #334155;
  border-radius: 10px;
  outline: none;
  font-family: inherit;
  font-size: 0.875rem;
  resize: vertical;
}

.form-group input::placeholder,
.form-group textarea::placeholder,
.comment-form textarea::placeholder {
  color: #64748b;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus,
.comment-form textarea:focus {
  border-color: #7c3aed;
  box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.18);
}

.form-group select option {
  background: #0f172a;
  color: #f8fafc;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1.25rem 1.75rem;
  background: #0f172a;
  border-top: 1px solid #1e293b;
  border-radius: 0 0 18px 18px;
}

.details-body {
  display: grid;
  grid-template-columns: 1fr 300px;
}

.details-main {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  padding: 1.5rem 1.75rem;
  border-right: 1px solid #1e293b;
}

.details-sidebar {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  padding: 1.5rem 1.25rem;
  background: #0f172a;
}

.details-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.detail-block-label {
  margin: 0 0 0.75rem;
  color: #94a3b8;
  font-size: 0.7rem;
  font-weight: 900;
  letter-spacing: 0.09em;
  text-transform: uppercase;
}

.detail-block-text {
  color: #cbd5e1;
  font-size: 0.875rem;
  line-height: 1.6;
}

.detail-meta-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem 1.5rem;
  padding: 1rem;
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 12px;
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.meta-label {
  color: #64748b;
  font-size: 0.72rem;
  font-weight: 800;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.meta-value {
  color: #f8fafc;
  font-size: 0.875rem;
  font-weight: 700;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 0.875rem;
  max-height: 280px;
  overflow-y: auto;
  margin-bottom: 0.875rem;
}

.comment-item {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
}

.comment-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  background: rgba(167, 139, 250, 0.16);
  color: #c4b5fd;
  border-radius: 999px;
  font-size: 0.7rem;
  font-weight: 900;
  flex-shrink: 0;
}

.comment-body {
  flex: 1;
  padding: 0.625rem 0.875rem;
  background: #0f172a;
  border: 1px solid #1e293b;
  border-radius: 12px;
}

.comment-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.35rem;
}

.comment-author {
  color: #f8fafc;
  font-size: 0.8rem;
  font-weight: 800;
}

.comment-type-badge {
  padding: 0.1rem 0.45rem;
  background: rgba(167, 139, 250, 0.16);
  color: #c4b5fd;
  border-radius: 6px;
  font-size: 0.65rem;
  font-weight: 800;
}

.comment-date {
  margin-left: auto;
  color: #64748b;
  font-size: 0.7rem;
}

.comment-text {
  color: #cbd5e1;
  font-size: 0.875rem;
  line-height: 1.5;
}

.empty-comments {
  color: #64748b;
  font-size: 0.82rem;
  font-style: italic;
}

.comment-form {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.attachment-list {
  display: flex;
  flex-direction: column;
  gap: 0.625rem;
  margin-bottom: 0.875rem;
}

.attachment-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.625rem 0.875rem;
  background: #111827;
  border: 1px solid #1e293b;
  border-radius: 10px;
}

.attachment-info {
  display: flex;
  align-items: center;
  gap: 0.625rem;
}

.attachment-name {
  color: #f8fafc;
  font-size: 0.82rem;
  font-weight: 800;
}

.attachment-size {
  color: #64748b;
  font-size: 0.72rem;
}

.attachment-actions {
  display: flex;
  gap: 0.375rem;
}

.link-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border: none;
  border-radius: 7px;
  cursor: pointer;
  font-size: 0.8rem;
  font-weight: 800;
}

.link-btn.blue {
  background: rgba(56, 189, 248, 0.14);
  color: #38bdf8;
}

.link-btn.red {
  background: rgba(251, 113, 133, 0.16);
  color: #fb7185;
}

.upload-zone {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.file-input {
  display: none;
}

.file-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 0.875rem;
  color: #94a3b8;
  border: 1.5px dashed #334155;
  border-radius: 10px;
  cursor: pointer;
  font-size: 0.8rem;
}

.file-label:hover {
  color: #c4b5fd;
  border-color: #7c3aed;
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.quick-btn {
  padding: 0.65rem 1rem;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  text-align: left;
  font-size: 0.85rem;
  font-weight: 800;
}

.quick-btn:hover {
  opacity: 0.9;
}

.quick-btn.indigo {
  background: rgba(167, 139, 250, 0.16);
  color: #c4b5fd;
}

.quick-btn.green {
  background: rgba(52, 211, 153, 0.14);
  color: #34d399;
}

.quick-btn.gray {
  background: rgba(148, 163, 184, 0.14);
  color: #cbd5e1;
}

@media (max-width: 1100px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .details-body {
    grid-template-columns: 1fr;
  }

  .details-main {
    border-right: none;
    border-bottom: 1px solid #1e293b;
  }
}

@media (max-width: 720px) {
  .incidents-page {
    padding: 1rem;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .filter-group {
    width: 100%;
    flex-direction: column;
  }

  .filter-select {
    width: 100%;
  }

  :deep(.dt-top),
  :deep(.dt-bottom) {
    align-items: stretch;
    flex-direction: column;
  }

  :deep(.dt-search input) {
    width: 100%;
    min-width: 0;
  }
}
</style>
