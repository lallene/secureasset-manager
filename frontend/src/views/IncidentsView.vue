<script setup>
import { onMounted, ref } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const incidents = ref([]);
const assets = ref([]);
const showModal = ref(false);
const isEditMode = ref(false);
const editingIncidentId = ref(null);
const error = ref("");
const role = localStorage.getItem("role");

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
const technicians = ref([]);

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

onMounted(() => {
  fetchIncidents();
  fetchTechnicians();
  fetchAssets();
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
            <th class="p-3 text-left">Asset</th>
            <th class="p-3 text-right">Actions</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="incident in incidents" :key="incident.ID" class="border-t">
            <td class="p-3">{{ incident.title }}</td>
            <td class="p-3">{{ incident.type }}</td>
            <td class="p-3">{{ incident.priority }}</td>
            <td class="p-3">{{ incident.status }}</td>
            <td class="p-3">
              {{ incident.asset?.name || "Non lié" }}
            </td>
            <td class="p-3 text-right">
             <button
                     v-if="role === 'Technician'"
                     @click="openEditModal(incident)"
                   >
                     Traiter
              </button>
              <button
                @click="openEditModal(incident)"
                class="text-blue-600 hover:underline mr-4"
              >
                Modifier
              </button>

              <button
                @click="deleteIncident(incident)"
                class="text-red-600 hover:underline"
              >
                Supprimer
              </button>
            </td>
          </tr>

          <tr v-if="incidents.length === 0">
            <td colspan="6" class="p-6 text-center text-gray-500">
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
            :value="technician.name"
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
</template>