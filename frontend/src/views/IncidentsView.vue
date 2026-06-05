<script setup>
import { onMounted, ref } from "vue";
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
const role = localStorage.getItem("role");

const highlightedIncidentId = ref(null);
const route = useRoute();

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
    const response = await api.get("/technicians", {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    technicians.value = response.data;
  } catch (err) {
    console.error(err);
  }
};

const fetchIncidents = async () => {
  try {
    const response = await api.get("/incidents", {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    incidents.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger les incidents";
  }
};

const fetchAssets = async () => {
  try {
    const response = await api.get("/assets", {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    assets.value = response.data;
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
    title: incident.title,
    description: incident.description,
    type: incident.type,
    priority: incident.priority,
    status: incident.status,
    asset_id: incident.asset_id,
    assigned_to: incident.assigned_to,
  };

  showModal.value = true;
};

const openDetailsModal = async (incident) => {
  selectedIncident.value = incident;
  comments.value = [];
  showDetailsModal.value = true;

  await fetchComments(incident.ID);
  await fetchAttachments(incident.ID);
};

const saveIncident = async () => {
  try {
    const payload = {
      ...form.value,
      asset_id: Number(form.value.asset_id),
    };

    if (isEditMode.value) {
      await api.put(`/incidents/${editingIncidentId.value}`, payload, {
        headers: {
          Authorization: `Bearer ${getToken()}`,
        },
      });
    } else {
      await api.post("/incidents", payload, {
        headers: {
          Authorization: `Bearer ${getToken()}`,
        },
      });
    }

    showModal.value = false;
    resetForm();
    fetchIncidents();
  } catch (err) {
    console.error(err);
    error.value = "Erreur lors de l'enregistrement de l'incident";
  }
};

const deleteIncident = async (incident) => {
  const confirmed = confirm(`Supprimer l'incident "${incident.title}" ?`);

  if (!confirmed) {
    return;
  }

  try {
    await api.delete(`/incidents/${incident.ID}`, {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    fetchIncidents();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de supprimer cet incident";
  }
};

const takeIncident = async (incident) => {
  try {
    await api.put(`/incidents/${incident.ID}/take`, {}, {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    fetchIncidents();
  } catch (err) {
    console.error(err);
    error.value = "Impossible d'assigner le ticket";
  }
};

const resolveIncident = async (incident) => {
  try {
    await api.put(`/incidents/${incident.ID}/resolve`, {}, {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    fetchIncidents();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de résoudre le ticket";
  }
};

const comments = ref([]);
const commentMessage = ref("");
const attachments = ref([]);
const selectedFile = ref(null);

const fetchComments = async (incidentId) => {
  try {
    const response = await api.get(`/incidents/${incidentId}/comments`, {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    comments.value = response.data || [];
  } catch (err) {
    console.error(err);
    comments.value = [];
  }
};

const addComment = async () => {
  if (!commentMessage.value.trim()) {
    return;
  }

  try {
    await api.post(
      `/incidents/${selectedIncident.value.ID}/comments`,
      {
        message: commentMessage.value,
      },
      {
        headers: {
          Authorization: `Bearer ${getToken()}`,
        },
      }
    );

    commentMessage.value = "";
    fetchComments(selectedIncident.value.ID);
  } catch (err) {
    console.error(err);
    error.value = "Impossible d'ajouter le commentaire";
  }
};

const fetchAttachments = async (incidentId) => {
  try {
    const response = await api.get(
      `/incidents/${incidentId}/attachments`,
      {
        headers: {
          Authorization: `Bearer ${getToken()}`,
        },
      }
    );

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
  if (!selectedFile.value) {
    return;
  }

  try {
    const formData = new FormData();

    formData.append("file", selectedFile.value);

    await api.post(
      `/incidents/${selectedIncident.value.ID}/attachments`,
      formData,
      {
        headers: {
          Authorization: `Bearer ${getToken()}`,
          "Content-Type": "multipart/form-data",
        },
      }
    );

    selectedFile.value = null;

    await fetchAttachments(selectedIncident.value.ID);
  } catch (err) {
    console.error(err);
  }
};

const downloadAttachment = async (attachment) => {
  try {
    const response = await api.get(
      `/attachments/${attachment.ID}/download`,
      {
        headers: {
          Authorization: `Bearer ${getToken()}`,
        },
        responseType: "blob",
      }
    );

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
  if (!confirm(`Supprimer ${attachment.file_name} ?`)) {
    return;
  }

  try {
    await api.delete(`/attachments/${attachment.ID}`, {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    await fetchAttachments(selectedIncident.value.ID);
  } catch (err) {
    console.error(err);
  }
};

const closeIncident = async (incident) => {
  try {
    await api.put(`/incidents/${incident.ID}/close`, {}, {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    fetchIncidents();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de clôturer le ticket";
  }
};

const getSlaStatus = (incident) => {
  if (incident.status === "Closed" || incident.status === "Resolved") {
    return {
      label: "Terminé",
      class: "text-slate-600",
    };
  }

  if (!incident.due_at) {
    return {
      label: "Non défini",
      class: "text-gray-500",
    };
  }

  const now = new Date();
  const dueDate = new Date(incident.due_at);
  const diffHours = (dueDate - now) / (1000 * 60 * 60);

  if (diffHours < 0) {
    return {
      label: "🔴 En retard",
      class: "text-red-600 font-semibold",
    };
  }

  if (diffHours <= 24) {
    return {
      label: "🟠 Échéance proche",
      class: "text-orange-600 font-semibold",
    };
  }

  return {
    label: "🟢 Dans les délais",
    class: "text-green-600 font-semibold",
  };
};

onMounted(async () => {
  await fetchIncidents();
  await fetchTechnicians();
  await fetchAssets();
  highlightedIncidentId.value = Number(route.query.incident);

  const incidentId = route.query.incident;

  if (incidentId) {
    const incident = incidents.value.find(
      (i) => i.ID === Number(incidentId)
    );

    if (incident) {
      openDetailsModal(incident);
    }
  }
});
</script>

<template>
  <DashboardLayout>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">Incidents</h1>

      <button
        v-if="role === 'Viewer'"
        @click="openCreateModal"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
      >
        Ajouter Incident
      </button>
    </div>

    <p v-if="error" class="mb-4 text-red-600">
      {{ error }}
    </p>

    <div class="bg-white rounded-xl shadow overflow-hidden">
      <table class="w-full">
        <thead class="bg-slate-100">
          <tr>
            <th class="p-3 text-left">Titre</th>
            <th class="p-3 text-left">Type</th>
            <th class="p-3 text-left">Priorité</th>
            <th class="p-3 text-left">Statut</th>
            <th class="p-3 text-left">SLA</th>
            <th class="p-3 text-left">Asset</th>
            <th class="p-3 text-right">Actions</th>
          </tr>
        </thead>

        <tbody>
          <tr
            v-for="incident in incidents"
            :key="incident.ID"
            class="border-t"
            :class="incident.ID === highlightedIncidentId ? 'bg-yellow-100' : ''"
          >
            <td class="p-3">{{ incident.title }}</td>
            <td class="p-3">{{ incident.type }}</td>
            <td class="p-3">{{ incident.priority }}</td>
            <td class="p-3">{{ incident.status }}</td>
            <td class="p-3">
              <span :class="getSlaStatus(incident).class">
                {{ getSlaStatus(incident).label }}
              </span>
            </td>
            <td class="p-3">
              {{ incident.asset?.name || "Non lié" }}
            </td>

            <td class="p-3 text-right">
              <button
                @click="openDetailsModal(incident)"
                class="text-slate-600 hover:underline mr-4"
              >
                Détails
              </button>

              <button
                v-if="(role === 'Admin' || role === 'Viewer') && incident.status !== 'Closed'"
                @click="openEditModal(incident)"
                class="text-blue-600 hover:underline mr-4"
              >
                Modifier
              </button>

             <button
               v-if="(role === 'Admin' || role === 'Viewer') && incident.status !== 'Closed'"
               @click="deleteIncident(incident)"
               class="text-red-600 hover:underline mr-4"
             >
               Supprimer
             </button>

              <button
                v-if="role === 'Technician' && incident.status !== 'Resolved' && incident.status !== 'Closed'"
                @click="takeIncident(incident)"
                class="text-indigo-600 hover:underline mr-4"
              >
                Assigner
              </button>

              <button
                v-if="role === 'Technician' && incident.status === 'In Progress'"
                @click="resolveIncident(incident)"
                class="text-green-600 hover:underline mr-4"
              >
                Résoudre
              </button>

              <button
                v-if="role === 'Technician' && incident.status === 'Resolved'"
                @click="closeIncident(incident)"
                class="text-orange-600 hover:underline"
              >
                Clôturer
              </button>
            </td>
          </tr>

          <tr v-if="incidents.length === 0">
            <td colspan="7" class="p-6 text-center text-gray-500">
              Aucun incident trouvé.
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </DashboardLayout>

  <div
    v-if="showModal"
    class="fixed inset-0 bg-black/50 flex items-center justify-center"
  >
    <div class="bg-white rounded-xl p-6 w-full max-w-lg">
      <h2 class="text-2xl font-bold mb-6">
        {{ isEditMode ? "Modifier Incident" : "Ajouter Incident" }}
      </h2>

      <div class="grid gap-4">
        <input
          v-model="form.title"
          placeholder="Titre"
          class="border p-3 rounded-lg"
        />

        <textarea
          v-model="form.description"
          placeholder="Description"
          class="border p-3 rounded-lg"
        ></textarea>

        <input
          v-model="form.type"
          placeholder="Type"
          class="border p-3 rounded-lg"
        />

        <select
          v-model="form.priority"
          class="border p-3 rounded-lg"
        >
          <option value="">Priorité</option>
          <option value="Low">Low</option>
          <option value="Medium">Medium</option>
          <option value="High">High</option>
          <option value="Critical">Critical</option>
        </select>

        <select
          v-if="role !== 'Viewer'"
          v-model="form.status"
          class="border p-3 rounded-lg"
        >
          <option value="">Statut</option>
          <option value="Open">Open</option>
          <option value="In Progress">In Progress</option>
          <option value="Resolved">Resolved</option>
          <option value="Closed">Closed</option>
        </select>

        <select
          v-model="form.asset_id"
          class="border p-3 rounded-lg"
        >
          <option value="">Asset lié</option>

          <option
            v-for="asset in assets"
            :key="asset.ID"
            :value="asset.ID"
          >
            {{ asset.name }} - {{ asset.ip_address }}
          </option>
        </select>

        <select
          v-model="form.assigned_to"
          class="border p-3 rounded-lg"
        >
          <option value="">Assigner à un technicien</option>

          <option
            v-for="technician in technicians"
            :key="technician.ID"
            :value="technician.email"
          >
            {{ technician.name }} - {{ technician.email }}
          </option>
        </select>
      </div>

      <div class="flex justify-end gap-4 mt-6">
        <button
          @click="showModal = false"
          class="px-4 py-2 border rounded-lg"
        >
          Annuler
        </button>

        <button
          @click="saveIncident"
          class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
        >
          {{ isEditMode ? "Enregistrer" : "Créer" }}
        </button>
      </div>
    </div>
  </div>

  <div
    v-if="showDetailsModal"
    class="fixed inset-0 bg-black/50 flex items-center justify-center"
  >
    <div class="bg-white rounded-xl p-6 w-full max-w-xl">
      <h2 class="text-2xl font-bold mb-4">Détails Incident</h2>

      <div class="space-y-2">
        <p><strong>Titre :</strong> {{ selectedIncident?.title }}</p>
        <p><strong>Description :</strong> {{ selectedIncident?.description }}</p>
        <p><strong>Type :</strong> {{ selectedIncident?.type }}</p>
        <p><strong>Priorité :</strong> {{ selectedIncident?.priority }}</p>
        <p><strong>Statut :</strong> {{ selectedIncident?.status }}</p>
        <p><strong>Asset :</strong> {{ selectedIncident?.asset?.name || "Non lié" }}</p>
        <p><strong>Assigné à :</strong> {{ selectedIncident?.assigned_to || "Non assigné" }}</p>
        <p><strong>Échéance :</strong> {{ selectedIncident?.due_at || "Non définie" }}</p>
      </div>
      <div class="mt-6 border-t pt-4">
        <h3 class="font-bold mb-3">Historique / Commentaires</h3>
        <div class="mt-6 border-t pt-4">
          <h3 class="font-bold mb-3">
            Pièces jointes
          </h3>

          <div class="space-y-2 mb-4">
            <div
              v-for="attachment in attachments"
              :key="attachment.ID"
              class="flex items-center justify-between bg-slate-100 p-3 rounded"
            >
              <div>
                📎 {{ attachment.file_name }}
                <span class="text-xs text-gray-500">
                  ({{ Math.round(attachment.file_size / 1024) }} KB)
                </span>
              </div>

              <div class="flex gap-3">
                <button
                  @click="downloadAttachment(attachment)"
                  class="text-blue-600 hover:underline"
                >
                  Télécharger
                </button>

                <button
                  v-if="selectedIncident?.status !== 'Closed'"
                  @click="deleteAttachment(attachment)"
                  class="text-red-600 hover:underline"
                >
                  Supprimer
                </button>
              </div>
            </div>

            <p
              v-if="attachments.length === 0"
              class="text-sm text-gray-500"
            >
              Aucune pièce jointe.
            </p>
          </div>

          <div
            v-if="selectedIncident?.status !== 'Closed'"
            class="flex gap-3 items-center"
          >
            <input
              type="file"
              @change="onFileChange"
              class="border p-2 rounded"
            />

            <button
              @click="uploadAttachment"
              class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded"
            >
              Upload
            </button>
          </div>
        </div>

        <div class="space-y-3 max-h-64 overflow-y-auto">
        <div
          v-for="comment in comments"
          :key="comment.ID"
          class="bg-slate-100 rounded-lg p-3"
        >
          <div class="text-sm text-slate-500 mb-1">
            {{ comment.user?.name || "Système" }} — {{ comment.type }}
          </div>

          <p>{{ comment.message }}</p>
        </div>

          <p v-if="comments.length === 0" class="text-sm text-gray-500">
            Aucun commentaire pour ce ticket.
          </p>
        </div>

        <div
          v-if="selectedIncident?.status !== 'Closed' && role !== 'Admin'"
          class="mt-4"
        >
          <textarea
            v-model="commentMessage"
            placeholder="Ajouter un commentaire..."
            class="border p-3 rounded-lg w-full"
          ></textarea>

          <button
            @click="addComment"
            class="mt-3 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
          >
            Ajouter commentaire
          </button>
        </div>
      </div>

      <div class="flex justify-end mt-6">
        <button
          @click="showDetailsModal = false"
          class="px-4 py-2 border rounded-lg"
        >
          Fermer
        </button>
      </div>
    </div>
  </div>
</template>