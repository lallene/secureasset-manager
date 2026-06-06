<script setup>
import { onMounted, ref } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const services = ref([]);
const showModal = ref(false);
const isEditMode = ref(false);
const editingServiceId = ref(null);
const error = ref("");
const isSubmitting = ref(false);

const form = ref({
  name: "",
  description: "",
});

const getToken = () => localStorage.getItem("token");

const fetchServices = async () => {
  try {
    const response = await api.get("/services", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    services.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger les services";
  }
};

const resetForm = () => {
  form.value = {
    name: "",
    description: "",
  };
  isEditMode.value = false;
  editingServiceId.value = null;
};

const openCreateModal = () => {
  resetForm();
  showModal.value = true;
};

const openEditModal = (service) => {
  isEditMode.value = true;
  editingServiceId.value = service.ID;

  form.value = {
    name: service.name,
    description: service.description,
  };

  showModal.value = true;
};

const saveService = async () => {
  try {
    isSubmitting.value = true;
    error.value = "";

    if (isEditMode.value) {
      await api.put(`/services/${editingServiceId.value}`, form.value, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    } else {
      await api.post("/services", form.value, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    }

    showModal.value = false;
    resetForm();
    await fetchServices();
  } catch (err) {
    console.error(err);
    error.value = "Impossible d'enregistrer le service";
  } finally {
    isSubmitting.value = false;
  }
};

const deleteService = async (service) => {
  if (!confirm(`Supprimer le service "${service.name}" ?`)) return;

  try {
    await api.delete(`/services/${service.ID}`, {
      headers: { Authorization: `Bearer ${getToken()}` },
    });

    await fetchServices();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de supprimer le service";
  }
};

onMounted(fetchServices);
</script>

<template>
  <DashboardLayout>
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-2xl font-semibold text-slate-900">Services</h1>
        <p class="text-sm text-slate-500 mt-1">
          Gérez les services responsables du traitement des tickets.
        </p>
      </div>

      <button
          @click="openCreateModal"
          class="bg-slate-900 hover:bg-slate-800 text-white px-4 py-2.5 rounded-xl text-sm font-medium"
      >
        Ajouter un service
      </button>
    </div>

    <div v-if="error" class="mb-6 p-4 bg-rose-50 border border-rose-100 text-rose-700 rounded-xl text-sm">
      {{ error }}
    </div>

    <div class="bg-white rounded-2xl border border-slate-100 shadow-sm overflow-hidden">
      <table class="w-full text-left">
        <thead class="bg-slate-50 border-b border-slate-100">
        <tr>
          <th class="p-4 text-xs uppercase text-slate-400">Nom</th>
          <th class="p-4 text-xs uppercase text-slate-400">Description</th>
          <th class="p-4 text-xs uppercase text-slate-400 text-right">Actions</th>
        </tr>
        </thead>

        <tbody class="divide-y divide-slate-100">
        <tr v-for="service in services" :key="service.ID" class="hover:bg-slate-50">
          <td class="p-4 font-medium text-slate-900">{{ service.name }}</td>
          <td class="p-4 text-slate-500">{{ service.description || "-" }}</td>
          <td class="p-4 text-right">
            <button @click="openEditModal(service)" class="text-blue-600 hover:underline mr-4">
              Modifier
            </button>
            <button @click="deleteService(service)" class="text-red-600 hover:underline">
              Supprimer
            </button>
          </td>
        </tr>

        <tr v-if="services.length === 0">
          <td colspan="3" class="p-10 text-center text-slate-400">
            Aucun service trouvé.
          </td>
        </tr>
        </tbody>
      </table>
    </div>
  </DashboardLayout>

  <div v-if="showModal" class="fixed inset-0 bg-slate-950/40 flex items-center justify-center p-4 z-50">
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-lg">
      <div class="p-6 border-b border-slate-100 flex justify-between">
        <h2 class="text-lg font-semibold">
          {{ isEditMode ? "Modifier le service" : "Ajouter un service" }}
        </h2>
        <button @click="showModal = false">✕</button>
      </div>

      <form @submit.prevent="saveService" class="p-6 space-y-4">
        <input v-model="form.name" required placeholder="Nom du service" class="input" />
        <textarea v-model="form.description" placeholder="Description" class="input"></textarea>

        <div class="flex justify-end gap-3 pt-4">
          <button type="button" @click="showModal = false" class="btn-secondary">
            Annuler
          </button>
          <button type="submit" :disabled="isSubmitting" class="btn-primary">
            {{ isEditMode ? "Enregistrer" : "Créer" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.input {
  width: 100%;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 10px 12px;
  font-size: 14px;
}
.btn-primary {
  background: #0f172a;
  color: white;
  padding: 9px 16px;
  border-radius: 12px;
}
.btn-secondary {
  border: 1px solid #e2e8f0;
  padding: 9px 16px;
  border-radius: 12px;
}
</style>