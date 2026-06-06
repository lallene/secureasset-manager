<script setup>
import { onMounted, ref, computed } from "vue";
import { useRoute } from "vue-router";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

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
const role = localStorage.getItem("role");

const highlightedIncidentId = ref(null);
const route = useRoute();

const searchQuery = ref("");
const filterStatus = ref("all");
const filterPriority = ref("all");

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
  Low:      { color: "#0284c7", bg: "#e0f2fe", dot: "#0ea5e9", label: "Faible" },
  Medium:   { color: "#7c3aed", bg: "#ede9fe", dot: "#8b5cf6", label: "Moyenne" },
  High:     { color: "#d97706", bg: "#fef3c7", dot: "#f59e0b", label: "Haute" },
  Critical: { color: "#dc2626", bg: "#fee2e2", dot: "#ef4444", label: "Critique" },
};

const statusConfig = {
  Open:         { color: "#0369a1", bg: "#e0f2fe", dot: "#0ea5e9",  label: "Ouvert" },
  "In Progress":{ color: "#7c3aed", bg: "#ede9fe", dot: "#8b5cf6",  label: "En cours" },
  Resolved:     { color: "#059669", bg: "#d1fae5", dot: "#10b981",  label: "Résolu" },
  Closed:       { color: "#6b7280", bg: "#f3f4f6", dot: "#9ca3af",  label: "Clôturé" },
};

const typeOptions = ["Panne", "Réseau", "Sécurité", "Matériel", "Logiciel", "Accès", "Autre"];

const getBadgeStyle = (map, key) => {
  const cfg = map[key];
  if (!cfg) return {};
  return { color: cfg.color, backgroundColor: cfg.bg, "--dot": cfg.dot };
};

const stats = computed(() => ({
  total: incidents.value.length,
  open: incidents.value.filter(i => i.status === "Open").length,
  inProgress: incidents.value.filter(i => i.status === "In Progress").length,
  critical: incidents.value.filter(i => i.priority === "Critical" && i.status !== "Closed").length,
  overdue: incidents.value.filter(i => {
    if (!i.due_at || i.status === "Closed" || i.status === "Resolved") return false;
    return new Date(i.due_at) < new Date();
  }).length,
}));

const filteredIncidents = computed(() => {
  let list = incidents.value;
  if (filterStatus.value !== "all") list = list.filter(i => i.status === filterStatus.value);
  if (filterPriority.value !== "all") list = list.filter(i => i.priority === filterPriority.value);
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    list = list.filter(i =>
      i.title?.toLowerCase().includes(q) ||
      i.type?.toLowerCase().includes(q) ||
      i.asset?.name?.toLowerCase().includes(q) ||
      i.assigned_to?.toLowerCase().includes(q)
    );
  }
  return list;
});

const resetForm = () => {
  form.value = { title: "", description: "", type: "", priority: "", status: "", asset_id: "", assigned_to: "" };
  isEditMode.value = false;
  editingIncidentId.value = null;
  error.value = "";
};

const fetchTechnicians = async () => {
  try {
    const r = await api.get("/technicians", { headers: { Authorization: `Bearer ${getToken()}` } });
    technicians.value = r.data;
  } catch {}
};

const fetchIncidents = async () => {
  isLoading.value = true;
  try {
    const r = await api.get("/incidents", { headers: { Authorization: `Bearer ${getToken()}` } });
    incidents.value = r.data;
  } catch {
    error.value = "Impossible de charger les incidents";
  } finally {
    isLoading.value = false;
  }
};

const fetchAssets = async () => {
  try {
    const r = await api.get("/assets", { headers: { Authorization: `Bearer ${getToken()}` } });
    assets.value = r.data;
  } catch {}
};

const openCreateModal = () => { resetForm(); form.value.status = "Open"; showModal.value = true; };

const openEditModal = (incident) => {
  isEditMode.value = true;
  editingIncidentId.value = incident.ID;
  form.value = { title: incident.title, description: incident.description, type: incident.type, priority: incident.priority, status: incident.status, asset_id: incident.asset_id, assigned_to: incident.assigned_to };
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
    const payload = { ...form.value, asset_id: Number(form.value.asset_id) };
    if (isEditMode.value) {
      await api.put(`/incidents/${editingIncidentId.value}`, payload, { headers: { Authorization: `Bearer ${getToken()}` } });
    } else {
      await api.post("/incidents", payload, { headers: { Authorization: `Bearer ${getToken()}` } });
    }
    showModal.value = false;
    resetForm();
    fetchIncidents();
  } catch {
    error.value = "Erreur lors de l'enregistrement de l'incident";
  }
};

const deleteIncident = async (incident) => {
  if (!confirm(`Supprimer l'incident « ${incident.title} » ?`)) return;
  try {
    await api.delete(`/incidents/${incident.ID}`, { headers: { Authorization: `Bearer ${getToken()}` } });
    fetchIncidents();
  } catch {
    error.value = "Impossible de supprimer cet incident";
  }
};

const takeIncident = async (incident) => {
  try {
    await api.put(`/incidents/${incident.ID}/take`, {}, { headers: { Authorization: `Bearer ${getToken()}` } });
    fetchIncidents();
  } catch { error.value = "Impossible d'assigner le ticket"; }
};

const resolveIncident = async (incident) => {
  try {
    await api.put(`/incidents/${incident.ID}/resolve`, {}, { headers: { Authorization: `Bearer ${getToken()}` } });
    fetchIncidents();
  } catch { error.value = "Impossible de résoudre le ticket"; }
};

const closeIncident = async (incident) => {
  try {
    await api.put(`/incidents/${incident.ID}/close`, {}, { headers: { Authorization: `Bearer ${getToken()}` } });
    fetchIncidents();
  } catch { error.value = "Impossible de clôturer le ticket"; }
};

const comments = ref([]);
const commentMessage = ref("");
const attachments = ref([]);
const selectedFile = ref(null);

const fetchComments = async (id) => {
  try {
    const r = await api.get(`/incidents/${id}/comments`, { headers: { Authorization: `Bearer ${getToken()}` } });
    comments.value = r.data || [];
  } catch { comments.value = []; }
};

const addComment = async () => {
  if (!commentMessage.value.trim()) return;
  try {
    await api.post(`/incidents/${selectedIncident.value.ID}/comments`, { message: commentMessage.value }, { headers: { Authorization: `Bearer ${getToken()}` } });
    commentMessage.value = "";
    fetchComments(selectedIncident.value.ID);
  } catch { error.value = "Impossible d'ajouter le commentaire"; }
};

const fetchAttachments = async (id) => {
  try {
    const r = await api.get(`/incidents/${id}/attachments`, { headers: { Authorization: `Bearer ${getToken()}` } });
    attachments.value = r.data || [];
  } catch { attachments.value = []; }
};

const onFileChange = (e) => { selectedFile.value = e.target.files[0]; };

const uploadAttachment = async () => {
  if (!selectedFile.value) return;
  try {
    const fd = new FormData();
    fd.append("file", selectedFile.value);
    await api.post(`/incidents/${selectedIncident.value.ID}/attachments`, fd, { headers: { Authorization: `Bearer ${getToken()}`, "Content-Type": "multipart/form-data" } });
    selectedFile.value = null;
    await fetchAttachments(selectedIncident.value.ID);
  } catch {}
};

const downloadAttachment = async (att) => {
  try {
    const r = await api.get(`/attachments/${att.ID}/download`, { headers: { Authorization: `Bearer ${getToken()}` }, responseType: "blob" });
    const url = window.URL.createObjectURL(new Blob([r.data]));
    const a = document.createElement("a");
    a.href = url; a.setAttribute("download", att.file_name);
    document.body.appendChild(a); a.click(); a.remove();
    window.URL.revokeObjectURL(url);
  } catch { error.value = "Impossible de télécharger la pièce jointe"; }
};

const deleteAttachment = async (att) => {
  if (!confirm(`Supprimer ${att.file_name} ?`)) return;
  try {
    await api.delete(`/attachments/${att.ID}`, { headers: { Authorization: `Bearer ${getToken()}` } });
    await fetchAttachments(selectedIncident.value.ID);
  } catch {}
};

const getSlaStatus = (incident) => {
  if (incident.status === "Closed" || incident.status === "Resolved")
    return { label: "Terminé", color: "#6b7280", bg: "#f3f4f6", icon: "✓" };
  if (!incident.due_at)
    return { label: "Non défini", color: "#94a3b8", bg: "#f8fafc", icon: "—" };
  const diffH = (new Date(incident.due_at) - new Date()) / 3600000;
  if (diffH < 0)   return { label: "En retard",      color: "#dc2626", bg: "#fee2e2", icon: "⚠" };
  if (diffH <= 24) return { label: "Échéance proche", color: "#d97706", bg: "#fef3c7", icon: "⏱" };
  return               { label: "Dans les délais",  color: "#059669", bg: "#d1fae5", icon: "✓" };
};

const formatDate = (d) => d ? new Date(d).toLocaleDateString("fr-FR", { day: "2-digit", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit" }) : "—";
const initials = (name) => name ? name.split(" ").map(p => p[0]).join("").toUpperCase().slice(0, 2) : "?";

onMounted(async () => {
  await fetchIncidents();
  await fetchTechnicians();
  await fetchAssets();
  highlightedIncidentId.value = Number(route.query.incident);
  const id = route.query.incident;
  if (id) {
    const inc = incidents.value.find(i => i.ID === Number(id));
    if (inc) openDetailsModal(inc);
  }
});
</script>

<template>
  <DashboardLayout>
    <div class="incidents-page">

      <!-- Header -->
      <div class="page-header">
        <div>
          <p class="page-eyebrow">Gestion ITSM</p>
          <h1 class="page-title">Incidents</h1>
        </div>
        <button v-if="role === 'Viewer'" class="btn-primary" @click="openCreateModal">
          <span>＋</span> Nouvel incident
        </button>
      </div>

      <!-- Stats -->
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon neutral">🎫</div>
          <div><div class="stat-value">{{ stats.total }}</div><div class="stat-label">Total</div></div>
        </div>
        <div class="stat-card">
          <div class="stat-icon blue">📂</div>
          <div><div class="stat-value blue">{{ stats.open }}</div><div class="stat-label">Ouverts</div></div>
        </div>
        <div class="stat-card">
          <div class="stat-icon purple">⚙</div>
          <div><div class="stat-value purple">{{ stats.inProgress }}</div><div class="stat-label">En cours</div></div>
        </div>
        <div class="stat-card">
          <div class="stat-icon red">🔴</div>
          <div><div class="stat-value red">{{ stats.critical }}</div><div class="stat-label">Critiques</div></div>
        </div>
        <div class="stat-card">
          <div class="stat-icon amber">⏰</div>
          <div><div class="stat-value amber">{{ stats.overdue }}</div><div class="stat-label">En retard</div></div>
        </div>
      </div>

      <!-- Toolbar -->
      <div class="toolbar">
        <div class="search-wrap">
          <span class="search-icon">🔍</span>
          <input v-model="searchQuery" type="text" placeholder="Rechercher titre, type, asset…" class="search-input" />
        </div>
        <div class="filter-group">
          <select v-model="filterStatus" class="filter-select">
            <option value="all">Tous les statuts</option>
            <option v-for="(v, k) in statusConfig" :key="k" :value="k">{{ v.label }}</option>
          </select>
          <select v-model="filterPriority" class="filter-select">
            <option value="all">Toutes priorités</option>
            <option v-for="(v, k) in priorityConfig" :key="k" :value="k">{{ v.label }}</option>
          </select>
        </div>
      </div>

      <div v-if="error" class="alert-error"><span>⚠</span> {{ error }}</div>

      <!-- Table -->
      <div class="table-card">
        <div v-if="isLoading" class="table-loading">
          <span class="spinner"></span> Chargement des incidents…
        </div>

        <table v-else class="inc-table">
          <thead>
            <tr>
              <th>Incident</th>
              <th>Priorité</th>
              <th>Statut</th>
              <th>SLA</th>
              <th>Asset lié</th>
              <th>Assigné à</th>
              <th>Créé le</th>
              <th class="col-actions"></th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="inc in filteredIncidents"
              :key="inc.ID"
              class="table-row"
              :class="{ highlighted: inc.ID === highlightedIncidentId }"
            >
              <!-- Incident title + type -->
              <td class="cell-incident">
                <div class="incident-title">{{ inc.title }}</div>
                <div class="incident-meta">
                  <span class="type-tag">{{ inc.type }}</span>
                  <span class="incident-id">#{{ inc.ID }}</span>
                </div>
              </td>

              <!-- Priority -->
              <td>
                <span class="badge" :style="getBadgeStyle(priorityConfig, inc.priority)">
                  <span class="badge-dot"></span>
                  {{ priorityConfig[inc.priority]?.label || inc.priority }}
                </span>
              </td>

              <!-- Status -->
              <td>
                <span class="badge" :style="getBadgeStyle(statusConfig, inc.status)">
                  <span class="badge-dot"></span>
                  {{ statusConfig[inc.status]?.label || inc.status }}
                </span>
              </td>

              <!-- SLA -->
              <td>
                <div class="sla-cell" :style="{ color: getSlaStatus(inc).color }">
                  <span class="sla-icon">{{ getSlaStatus(inc).icon }}</span>
                  <span class="sla-label">{{ getSlaStatus(inc).label }}</span>
                </div>
                <div v-if="inc.due_at" class="sla-date">{{ formatDate(inc.due_at) }}</div>
              </td>

              <!-- Asset -->
              <td>
                <span v-if="inc.asset?.name" class="asset-chip">🖥 {{ inc.asset.name }}</span>
                <span v-else class="cell-muted">Non lié</span>
              </td>

              <!-- Assigned -->
              <td>
                <div v-if="inc.assigned_to" class="assigned-user">
                  <span class="user-avatar">{{ initials(inc.assigned_to) }}</span>
                  <span class="user-name">{{ inc.assigned_to }}</span>
                </div>
                <span v-else class="cell-muted unassigned">Non assigné</span>
              </td>

              <!-- Date -->
              <td class="cell-date">{{ formatDate(inc.CreatedAt || inc.created_at) }}</td>

              <!-- Actions -->
              <td class="cell-actions">
                <button class="action-btn" title="Voir les détails" @click="openDetailsModal(inc)">🔍</button>

                <button
                  v-if="(role === 'Admin' || role === 'Viewer') && inc.status !== 'Closed'"
                  class="action-btn edit"
                  title="Modifier"
                  @click="openEditModal(inc)"
                >✏</button>

                <button
                  v-if="(role === 'Admin' || role === 'Viewer') && inc.status !== 'Closed'"
                  class="action-btn del"
                  title="Supprimer"
                  @click="deleteIncident(inc)"
                >🗑</button>

                <button
                  v-if="role === 'Technician' && inc.status !== 'Resolved' && inc.status !== 'Closed'"
                  class="action-btn take"
                  title="Prendre en charge"
                  @click="takeIncident(inc)"
                >⚡</button>

                <button
                  v-if="role === 'Technician' && inc.status === 'In Progress'"
                  class="action-btn resolve"
                  title="Résoudre"
                  @click="resolveIncident(inc)"
                >✅</button>

                <button
                  v-if="role === 'Technician' && inc.status === 'Resolved'"
                  class="action-btn close-btn"
                  title="Clôturer"
                  @click="closeIncident(inc)"
                >🔒</button>
              </td>
            </tr>

            <tr v-if="filteredIncidents.length === 0 && !isLoading">
              <td colspan="8" class="empty-state">
                <span class="empty-icon">🎫</span>
                <p>Aucun incident trouvé</p>
                <small>Modifiez vos filtres ou créez un nouvel incident.</small>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="filteredIncidents.length > 0" class="table-footer">
          {{ filteredIncidents.length }} incident{{ filteredIncidents.length > 1 ? 's' : '' }}
          <span v-if="incidents.length !== filteredIncidents.length"> sur {{ incidents.length }} au total</span>
        </div>
      </div>
    </div>
  </DashboardLayout>

  <!-- ─── Création / Édition Modal ─────────────────────────────── -->
  <Teleport to="body">
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <div class="modal-header">
          <div>
            <p class="modal-eyebrow">{{ isEditMode ? 'Modification' : 'Création' }}</p>
            <h2 class="modal-title">{{ isEditMode ? 'Modifier l\'incident' : 'Nouvel incident' }}</h2>
          </div>
          <button class="modal-close" @click="showModal = false">✕</button>
        </div>

        <div v-if="error" class="alert-error modal-alert"><span>⚠</span> {{ error }}</div>

        <div class="modal-body">
          <!-- Section principale -->
          <div class="form-section">
            <p class="form-section-label">Identification</p>
            <div class="form-group full">
              <label>Titre <span class="required">*</span></label>
              <input v-model="form.title" type="text" placeholder="Résumé concis de l'incident" />
            </div>
            <div class="form-group full">
              <label>Description</label>
              <textarea v-model="form.description" rows="3" placeholder="Contexte, symptômes, actions déjà tentées…"></textarea>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label>Type</label>
                <select v-model="form.type">
                  <option value="" disabled>Sélectionner…</option>
                  <option v-for="t in typeOptions" :key="t" :value="t">{{ t }}</option>
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

          <!-- Section affectation -->
          <div class="form-section">
            <p class="form-section-label">Affectation</p>
            <div class="form-row">
              <div class="form-group">
                <label>Statut</label>
                <select v-if="role !== 'Viewer'" v-model="form.status">
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
                  <option v-for="t in technicians" :key="t.ID" :value="t.email">
                    {{ t.name }} — {{ t.email }}
                  </option>
                </select>
              </div>
            </div>
            <div class="form-group full">
              <label>Asset concerné</label>
              <select v-model="form.asset_id">
                <option value="">Aucun asset lié</option>
                <option v-for="a in assets" :key="a.ID" :value="a.ID">{{ a.name }} — {{ a.ip_address }}</option>
              </select>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn-ghost" @click="showModal = false">Annuler</button>
          <button class="btn-primary" @click="saveIncident">
            {{ isEditMode ? 'Enregistrer les modifications' : 'Créer l\'incident' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- ─── Détails Modal ──────────────────────────────────────────── -->
  <Teleport to="body">
    <div v-if="showDetailsModal" class="modal-overlay" @click.self="showDetailsModal = false">
      <div class="modal modal-lg">
        <div class="modal-header">
          <div style="flex:1; min-width:0;">
            <p class="modal-eyebrow">Incident #{{ selectedIncident?.ID }}</p>
            <h2 class="modal-title truncate">{{ selectedIncident?.title }}</h2>
          </div>
          <button class="modal-close" @click="showDetailsModal = false">✕</button>
        </div>

        <div class="details-body">
          <!-- Left column: infos + commentaires -->
          <div class="details-main">

            <!-- Badges -->
            <div class="details-badges">
              <span class="badge" :style="getBadgeStyle(priorityConfig, selectedIncident?.priority)">
                <span class="badge-dot"></span>
                {{ priorityConfig[selectedIncident?.priority]?.label || selectedIncident?.priority }}
              </span>
              <span class="badge" :style="getBadgeStyle(statusConfig, selectedIncident?.status)">
                <span class="badge-dot"></span>
                {{ statusConfig[selectedIncident?.status]?.label || selectedIncident?.status }}
              </span>
              <span class="badge sla-badge" :style="{ color: getSlaStatus(selectedIncident).color, backgroundColor: getSlaStatus(selectedIncident).bg }">
                {{ getSlaStatus(selectedIncident).icon }} {{ getSlaStatus(selectedIncident).label }}
              </span>
            </div>

            <!-- Description -->
            <div class="detail-block">
              <p class="detail-block-label">Description</p>
              <p class="detail-block-text">{{ selectedIncident?.description || 'Aucune description fournie.' }}</p>
            </div>

            <!-- Meta grid -->
            <div class="detail-meta-grid">
              <div class="meta-item">
                <span class="meta-label">Type</span>
                <span class="meta-value">{{ selectedIncident?.type || '—' }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">Asset</span>
                <span class="meta-value">{{ selectedIncident?.asset?.name || 'Non lié' }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">Assigné à</span>
                <span class="meta-value">{{ selectedIncident?.assigned_to || 'Non assigné' }}</span>
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

            <!-- Commentaires -->
            <div class="detail-block">
              <p class="detail-block-label">Historique & commentaires ({{ comments.length }})</p>
              <div class="comments-list">
                <div v-for="c in comments" :key="c.ID" class="comment-item">
                  <div class="comment-avatar">{{ initials(c.user?.name || 'Sys') }}</div>
                  <div class="comment-body">
                    <div class="comment-meta">
                      <span class="comment-author">{{ c.user?.name || 'Système' }}</span>
                      <span class="comment-type-badge">{{ c.type }}</span>
                      <span class="comment-date">{{ formatDate(c.CreatedAt) }}</span>
                    </div>
                    <p class="comment-text">{{ c.message }}</p>
                  </div>
                </div>
                <p v-if="comments.length === 0" class="empty-comments">Aucun commentaire pour ce ticket.</p>
              </div>

              <div v-if="selectedIncident?.status !== 'Closed' && role !== 'Admin'" class="comment-form">
                <textarea v-model="commentMessage" rows="2" placeholder="Ajouter un commentaire…"></textarea>
                <button class="btn-primary btn-sm" @click="addComment">Publier</button>
              </div>
            </div>
          </div>

          <!-- Right column: attachments -->
          <div class="details-sidebar">
            <div class="detail-block">
              <p class="detail-block-label">Pièces jointes ({{ attachments.length }})</p>

              <div class="attachment-list">
                <div v-for="att in attachments" :key="att.ID" class="attachment-item">
                  <div class="attachment-info">
                    <span class="attachment-icon">📎</span>
                    <div>
                      <div class="attachment-name">{{ att.file_name }}</div>
                      <div class="attachment-size">{{ Math.round(att.file_size / 1024) }} KB</div>
                    </div>
                  </div>
                  <div class="attachment-actions">
                    <button class="link-btn blue" @click="downloadAttachment(att)">↓</button>
                    <button
                      v-if="selectedIncident?.status !== 'Closed'"
                      class="link-btn red"
                      @click="deleteAttachment(att)"
                    >✕</button>
                  </div>
                </div>
                <p v-if="attachments.length === 0" class="empty-comments">Aucune pièce jointe.</p>
              </div>

              <div v-if="selectedIncident?.status !== 'Closed'" class="upload-zone">
                <input type="file" id="file-upload" class="file-input" @change="onFileChange" />
                <label for="file-upload" class="file-label">
                  <span>📁</span>
                  <span>{{ selectedFile ? selectedFile.name : 'Choisir un fichier…' }}</span>
                </label>
                <button v-if="selectedFile" class="btn-primary btn-sm" @click="uploadAttachment">Uploader</button>
              </div>
            </div>

            <!-- Quick actions -->
            <div class="detail-block" v-if="role === 'Technician'">
              <p class="detail-block-label">Actions rapides</p>
              <div class="quick-actions">
                <button
                  v-if="selectedIncident?.status !== 'Resolved' && selectedIncident?.status !== 'Closed'"
                  class="quick-btn indigo"
                  @click="takeIncident(selectedIncident); showDetailsModal = false;"
                >⚡ Prendre en charge</button>
                <button
                  v-if="selectedIncident?.status === 'In Progress'"
                  class="quick-btn green"
                  @click="resolveIncident(selectedIncident); showDetailsModal = false;"
                >✅ Marquer résolu</button>
                <button
                  v-if="selectedIncident?.status === 'Resolved'"
                  class="quick-btn gray"
                  @click="closeIncident(selectedIncident); showDetailsModal = false;"
                >🔒 Clôturer</button>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn-ghost" @click="showDetailsModal = false">Fermer</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
/* ─── Base ──────────────────────────────────────────────── */
.incidents-page {
  padding: 2rem;
  max-width: 1280px;
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
  font-size: 0.72rem;
  font-weight: 600;
  letter-spacing: 0.09em;
  text-transform: uppercase;
  color: #7c3aed;
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
.stats-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 1rem;
  margin-bottom: 1.5rem;
}
.stat-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 1rem 1.125rem;
  display: flex;
  align-items: center;
  gap: 0.875rem;
}
.stat-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.125rem;
  flex-shrink: 0;
}
.stat-icon.neutral { background: #f1f5f9; }
.stat-icon.blue    { background: #e0f2fe; }
.stat-icon.purple  { background: #ede9fe; }
.stat-icon.red     { background: #fee2e2; }
.stat-icon.amber   { background: #fef3c7; }
.stat-value { font-size: 1.5rem; font-weight: 700; color: #0f172a; line-height: 1; }
.stat-value.blue   { color: #0284c7; }
.stat-value.purple { color: #7c3aed; }
.stat-value.red    { color: #dc2626; }
.stat-value.amber  { color: #d97706; }
.stat-label { font-size: 0.78rem; color: #64748b; font-weight: 500; margin-top: 0.1rem; }

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
.search-input:focus { border-color: #a78bfa; box-shadow: 0 0 0 3px rgba(124, 58, 237, 0.1); }
.filter-group { display: flex; gap: 0.5rem; }
.filter-select {
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 0.6rem 2rem 0.6rem 0.875rem;
  font-size: 0.85rem;
  color: #374151;
  background: #fff;
  outline: none;
  cursor: pointer;
  -webkit-appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%2394a3b8' d='M6 8L1 3h10z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 0.75rem center;
}
.filter-select:focus { border-color: #a78bfa; }

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
.modal-alert { border-radius: 0; border-left: none; border-right: none; border-top: none; margin-bottom: 0; }

/* ─── Table ─────────────────────────────────────────────── */
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
  border-top-color: #7c3aed;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  display: inline-block;
}
@keyframes spin { to { transform: rotate(360deg); } }

.inc-table { width: 100%; border-collapse: collapse; font-size: 0.875rem; }
.inc-table thead tr {
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}
.inc-table th {
  padding: 0.75rem 1rem;
  text-align: left;
  font-size: 0.72rem;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  white-space: nowrap;
}
.col-actions { width: 120px; }
.table-row { border-bottom: 1px solid #f1f5f9; transition: background 0.1s; }
.table-row:last-child { border-bottom: none; }
.table-row:hover { background: #fdfaff; }
.table-row.highlighted { background: #fefce8; }
.inc-table td { padding: 0.875rem 1rem; color: #1e293b; vertical-align: middle; }

/* ─── Cell variants ─────────────────────────────────────── */
.cell-incident .incident-title { font-weight: 600; color: #0f172a; margin-bottom: 0.25rem; }
.cell-incident .incident-meta { display: flex; align-items: center; gap: 0.5rem; }
.type-tag {
  background: #f1f5f9;
  color: #475569;
  font-size: 0.7rem;
  font-weight: 500;
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
}
.incident-id { font-size: 0.72rem; color: #94a3b8; }

.badge {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.25rem 0.65rem;
  border-radius: 999px;
  white-space: nowrap;
}
.badge-dot { width: 6px; height: 6px; background: var(--dot, currentColor); border-radius: 50%; flex-shrink: 0; }

.sla-cell { display: flex; align-items: center; gap: 0.35rem; font-size: 0.8rem; font-weight: 600; }
.sla-date { font-size: 0.7rem; color: #94a3b8; margin-top: 0.15rem; }

.asset-chip { font-size: 0.8rem; color: #334155; }
.cell-muted { color: #94a3b8; font-size: 0.85rem; }
.unassigned { font-style: italic; }

.assigned-user { display: flex; align-items: center; gap: 0.5rem; }
.user-avatar {
  width: 26px; height: 26px;
  background: #ede9fe; color: #7c3aed;
  border-radius: 50%; font-size: 0.7rem; font-weight: 700;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.user-name { font-size: 0.82rem; color: #374151; }
.cell-date { font-size: 0.78rem; color: #64748b; white-space: nowrap; }

.cell-actions { display: flex; gap: 0.2rem; justify-content: flex-end; flex-wrap: wrap; }
.action-btn {
  border: none; background: transparent; padding: 0.35rem; border-radius: 7px;
  cursor: pointer; font-size: 0.9rem; opacity: 0.5; transition: opacity 0.15s, background 0.15s;
}
.action-btn:hover { opacity: 1; }
.action-btn.edit:hover   { background: #eff6ff; }
.action-btn.del:hover    { background: #fef2f2; }
.action-btn.take:hover   { background: #ede9fe; }
.action-btn.resolve:hover{ background: #d1fae5; }
.action-btn.close-btn:hover { background: #f3f4f6; }

.empty-state { text-align: center; padding: 3.5rem 1rem !important; }
.empty-icon { font-size: 2.5rem; display: block; margin-bottom: 0.75rem; }
.empty-state p { font-size: 1rem; font-weight: 500; color: #334155; margin: 0 0 0.25rem; }
.empty-state small { font-size: 0.8rem; color: #94a3b8; }
.table-footer { padding: 0.75rem 1rem; border-top: 1px solid #f1f5f9; font-size: 0.8rem; color: #94a3b8; background: #fafbff; }

/* ─── Buttons ───────────────────────────────────────────── */
.btn-primary {
  display: inline-flex; align-items: center; gap: 0.5rem;
  background: #7c3aed; color: #fff; border: none;
  padding: 0.625rem 1.125rem; border-radius: 10px;
  font-size: 0.875rem; font-weight: 600; cursor: pointer;
  transition: background 0.15s, box-shadow 0.15s, transform 0.1s;
}
.btn-primary:hover { background: #6d28d9; box-shadow: 0 4px 12px rgba(124,58,237,0.3); }
.btn-primary:active { transform: scale(0.98); }
.btn-sm { padding: 0.45rem 0.875rem; font-size: 0.8rem; }
.btn-ghost {
  background: transparent; border: 1px solid #e2e8f0; color: #374151;
  padding: 0.625rem 1.125rem; border-radius: 10px;
  font-size: 0.875rem; font-weight: 500; cursor: pointer;
  transition: background 0.15s;
}
.btn-ghost:hover { background: #f8fafc; }

/* ─── Modal ─────────────────────────────────────────────── */
.modal-overlay {
  position: fixed; inset: 0;
  background: rgba(15,23,42,0.45);
  display: flex; align-items: center; justify-content: center;
  z-index: 9999; backdrop-filter: blur(2px);
}
.modal {
  background: #fff; border-radius: 18px;
  width: 100%; max-width: 580px; max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 24px 64px rgba(0,0,0,0.15);
}
.modal-lg { max-width: 880px; }

.modal-header {
  display: flex; align-items: flex-start; justify-content: space-between;
  padding: 1.5rem 1.75rem 1.25rem;
  border-bottom: 1px solid #f1f5f9;
}
.modal-eyebrow { font-size: 0.7rem; font-weight: 600; letter-spacing: 0.1em; text-transform: uppercase; color: #7c3aed; margin: 0 0 0.2rem; }
.modal-title { font-size: 1.2rem; font-weight: 700; color: #0f172a; margin: 0; }
.truncate { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.modal-close {
  border: none; background: #f1f5f9; color: #64748b;
  width: 30px; height: 30px; border-radius: 8px;
  cursor: pointer; font-size: 0.875rem; display: flex; align-items: center; justify-content: center;
  transition: background 0.15s; flex-shrink: 0;
}
.modal-close:hover { background: #e2e8f0; color: #0f172a; }

/* ─── Form modal ────────────────────────────────────────── */
.modal-body { padding: 1.25rem 1.75rem; display: flex; flex-direction: column; gap: 1.25rem; }
.form-section-label {
  font-size: 0.72rem; font-weight: 700; letter-spacing: 0.08em;
  text-transform: uppercase; color: #94a3b8; margin: 0 0 0.875rem;
}
.form-section { display: flex; flex-direction: column; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.form-group { display: flex; flex-direction: column; gap: 0.375rem; margin-bottom: 0.875rem; }
.form-group.full { grid-column: 1 / -1; }
.form-group:last-child { margin-bottom: 0; }
.form-group label { font-size: 0.8rem; font-weight: 600; color: #374151; }
.required { color: #ef4444; }
.form-group input,
.form-group select,
.form-group textarea {
  border: 1px solid #e2e8f0; border-radius: 9px;
  padding: 0.6rem 0.875rem; font-size: 0.875rem;
  color: #0f172a; background: #fff; outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
  font-family: inherit; resize: vertical;
}
.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus { border-color: #a78bfa; box-shadow: 0 0 0 3px rgba(124,58,237,0.1); }
.form-group input::placeholder,
.form-group textarea::placeholder { color: #c0ccda; }

.modal-footer {
  display: flex; justify-content: flex-end; gap: 0.75rem;
  padding: 1.25rem 1.75rem; border-top: 1px solid #f1f5f9;
  background: #fafbff; border-radius: 0 0 18px 18px;
}

/* ─── Details modal ─────────────────────────────────────── */
.details-body {
  display: grid;
  grid-template-columns: 1fr 280px;
  gap: 0;
}
.details-main {
  padding: 1.5rem 1.75rem;
  border-right: 1px solid #f1f5f9;
  display: flex; flex-direction: column; gap: 1.25rem;
}
.details-sidebar {
  padding: 1.5rem 1.25rem;
  background: #fafbff;
  display: flex; flex-direction: column; gap: 1.25rem;
}
.details-badges { display: flex; flex-wrap: wrap; gap: 0.5rem; }

.detail-block {}
.detail-block-label {
  font-size: 0.7rem; font-weight: 700;
  letter-spacing: 0.09em; text-transform: uppercase;
  color: #94a3b8; margin: 0 0 0.75rem;
}
.detail-block-text { font-size: 0.875rem; color: #374151; line-height: 1.6; }

.detail-meta-grid {
  display: grid; grid-template-columns: 1fr 1fr;
  gap: 0.75rem 1.5rem;
  background: #f8fafc; border-radius: 10px; padding: 1rem;
}
.meta-item { display: flex; flex-direction: column; gap: 0.2rem; }
.meta-label { font-size: 0.72rem; font-weight: 600; color: #94a3b8; text-transform: uppercase; letter-spacing: 0.05em; }
.meta-value { font-size: 0.875rem; color: #0f172a; font-weight: 500; }

/* ─── Comments ──────────────────────────────────────────── */
.comments-list { display: flex; flex-direction: column; gap: 0.875rem; max-height: 280px; overflow-y: auto; margin-bottom: 0.875rem; }
.comment-item { display: flex; gap: 0.75rem; align-items: flex-start; }
.comment-avatar {
  width: 30px; height: 30px; border-radius: 50%;
  background: #ede9fe; color: #7c3aed;
  font-size: 0.7rem; font-weight: 700;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.comment-body { flex: 1; background: #f8fafc; border-radius: 10px; padding: 0.625rem 0.875rem; }
.comment-meta { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.35rem; flex-wrap: wrap; }
.comment-author { font-size: 0.8rem; font-weight: 600; color: #0f172a; }
.comment-type-badge {
  font-size: 0.65rem; font-weight: 600; background: #ede9fe;
  color: #7c3aed; padding: 0.1rem 0.45rem; border-radius: 4px;
}
.comment-date { font-size: 0.7rem; color: #94a3b8; margin-left: auto; }
.comment-text { font-size: 0.875rem; color: #374151; line-height: 1.5; }
.empty-comments { font-size: 0.82rem; color: #94a3b8; font-style: italic; }
.comment-form { display: flex; flex-direction: column; gap: 0.5rem; }
.comment-form textarea {
  border: 1px solid #e2e8f0; border-radius: 9px; padding: 0.6rem 0.875rem;
  font-size: 0.875rem; color: #0f172a; outline: none; resize: vertical;
  font-family: inherit; transition: border-color 0.15s;
}
.comment-form textarea:focus { border-color: #a78bfa; box-shadow: 0 0 0 3px rgba(124,58,237,0.1); }

/* ─── Attachments ───────────────────────────────────────── */
.attachment-list { display: flex; flex-direction: column; gap: 0.625rem; margin-bottom: 0.875rem; }
.attachment-item {
  display: flex; align-items: center; justify-content: space-between;
  background: #fff; border: 1px solid #e2e8f0; border-radius: 9px;
  padding: 0.625rem 0.875rem;
}
.attachment-info { display: flex; align-items: center; gap: 0.625rem; }
.attachment-icon { font-size: 1.1rem; }
.attachment-name { font-size: 0.82rem; font-weight: 500; color: #0f172a; }
.attachment-size { font-size: 0.72rem; color: #94a3b8; }
.attachment-actions { display: flex; gap: 0.375rem; }
.link-btn {
  width: 26px; height: 26px; border: none; border-radius: 6px;
  cursor: pointer; font-size: 0.8rem; font-weight: 600;
  display: flex; align-items: center; justify-content: center;
  transition: background 0.15s;
}
.link-btn.blue { background: #e0f2fe; color: #0284c7; }
.link-btn.blue:hover { background: #bae6fd; }
.link-btn.red { background: #fee2e2; color: #dc2626; }
.link-btn.red:hover { background: #fecaca; }

.upload-zone { display: flex; flex-direction: column; gap: 0.5rem; }
.file-input { display: none; }
.file-label {
  display: flex; align-items: center; gap: 0.5rem;
  border: 1.5px dashed #e2e8f0; border-radius: 9px;
  padding: 0.625rem 0.875rem; font-size: 0.8rem; color: #64748b;
  cursor: pointer; transition: border-color 0.15s;
}
.file-label:hover { border-color: #a78bfa; color: #7c3aed; }

/* ─── Quick actions ─────────────────────────────────────── */
.quick-actions { display: flex; flex-direction: column; gap: 0.5rem; }
.quick-btn {
  border: none; border-radius: 9px; padding: 0.625rem 1rem;
  font-size: 0.85rem; font-weight: 600; cursor: pointer; text-align: left;
  transition: opacity 0.15s, transform 0.1s;
}
.quick-btn:hover { opacity: 0.85; transform: translateX(2px); }
.quick-btn.indigo { background: #ede9fe; color: #6d28d9; }
.quick-btn.green  { background: #d1fae5; color: #065f46; }
.quick-btn.gray   { background: #f3f4f6; color: #374151; }
.sla-badge { font-weight: 600; }
</style>