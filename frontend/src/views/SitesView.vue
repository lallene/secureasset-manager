<script setup>
import { onMounted, ref } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const sites = ref([]);
const showModal = ref(false);
const isEditMode = ref(false);
const editingSiteId = ref(null);
const error = ref("");
const isSubmitting = ref(false);

const form = ref({
  name: "",
  city: "",
  country: "France",
  description: "",
});

const getToken = () => localStorage.getItem("token");

const fetchSites = async () => {
  try {
    const response = await api.get("/sites", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    sites.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger les sites";
  }
};

const resetForm = () => {
  form.value = {
    name: "",
    city: "",
    country: "France",
    description: "",
  };
  isEditMode.value = false;
  editingSiteId.value = null;
};

const openCreateModal = () => {
  resetForm();
  showModal.value = true;
};

const openEditModal = (site) => {
  isEditMode.value = true;
  editingSiteId.value = site.ID;

  form.value = {
    name: site.name,
    city: site.city,
    country: site.country,
    description: site.description,
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
    error.value = "Impossible d'enregistrer le site";
  } finally {
    isSubmitting.value = false;
  }
};

const deleteSite = async (site) => {
  if (!confirm(`Supprimer le site "${site.name}" ?`)) return;

  try {
    await api.delete(`/sites/${site.ID}`, {
      headers: { Authorization: `Bearer ${getToken()}` },
    });

    await fetchSites();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de supprimer le site";
  }
};

onMounted(fetchSites);
</script>

<template>
  <DashboardLayout>
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-2xl font-semibold text-slate-900">Sites</h1>
        <p class="text-sm text-slate-500 mt-1">
          Gérez les implantations de votre organisation.
        </p>
      </div>

      <button
          @click="openCreateModal"
          class="bg-slate-900 hover:bg-slate-800 text-white px-4 py-2.5 rounded-xl text-sm font-medium"
      >
        Ajouter un site
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
          <th class="p-4 text-xs uppercase text-slate-400">Ville</th>
          <th class="p-4 text-xs uppercase text-slate-400">Pays</th>
          <th class="p-4 text-xs uppercase text-slate-400">Description</th>
          <th class="p-4 text-xs uppercase text-slate-400 text-right">Actions</th>
        </tr>
        </thead>

        <tbody class="divide-y divide-slate-100">
        <tr v-for="site in sites" :key="site.ID" class="hover:bg-slate-50">
          <td class="p-4 font-medium text-slate-900">{{ site.name }}</td>
          <td class="p-4 text-slate-600">{{ site.city || "-" }}</td>
          <td class="p-4 text-slate-600">{{ site.country || "-" }}</td>
          <td class="p-4 text-slate-500">{{ site.description || "-" }}</td>
          <td class="p-4 text-right">
            <button @click="openEditModal(site)" class="text-blue-600 hover:underline mr-4">
              Modifier
            </button>
            <button @click="deleteSite(site)" class="text-red-600 hover:underline">
              Supprimer
            </button>
          </td>
        </tr>

        <tr v-if="sites.length === 0">
          <td colspan="5" class="p-10 text-center text-slate-400">
            Aucun site trouvé.
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
          {{ isEditMode ? "Modifier le site" : "Ajouter un site" }}
        </h2>
        <button @click="showModal = false">✕</button>
      </div>

      <form @submit.prevent="saveSite" class="p-6 space-y-4">
        <input v-model="form.name" required placeholder="Nom du site" class="input" />
        <input v-model="form.city" placeholder="Ville" class="input" />
        <input v-model="form.country" placeholder="Pays" class="input" />
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