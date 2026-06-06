<script setup>
import { onMounted, ref, computed } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

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

const criticalityConfig = {
  Low: "bg-slate-100 text-slate-700 border-slate-200",
  Medium: "bg-blue-50 text-blue-700 border-blue-100",
  High: "bg-amber-50 text-amber-700 border-amber-100",
  Critical: "bg-red-50 text-red-700 border-red-100",
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
    error.value = "Impossible de charger les applications";
  }
};

const fetchSites = async () => {
  try {
    const response = await api.get("/sites", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });

    sites.value = response.data;
  } catch (err) {
    console.error(err);
  }
};

const fetchUsers = async () => {
  try {
    const response = await api.get("/users", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });

    users.value = response.data;
  } catch (err) {
    console.error(err);
  }
};

const resetForm = () => {
  form.value = {
    name: "",
    description: "",
    version: "",
    criticality: "Medium",
    status: "Active",
    site_id: "",
    owner_id: "",
  };

  isEditMode.value = false;
  editingApplicationId.value = null;
  error.value = "";
};

const openCreateModal = () => {
  resetForm();
  showModal.value = true;
};

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
    console.error(err);
    error.value = "Impossible d'enregistrer l'application";
  } finally {
    isSubmitting.value = false;
  }
};

const deleteApplication = async (app) => {
  if (!confirm(`Supprimer l'application "${app.name}" ?`)) return;

  try {
    await api.delete(`/cmdb/applications/${app.ID}`, {
      headers: { Authorization: `Bearer ${getToken()}` },
    });

    await fetchApplications();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de supprimer l'application";
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
    <div class="p-8">
      <div class="flex items-center justify-between mb-8">
        <div>
          <p class="text-xs font-semibold uppercase tracking-wider text-indigo-600">
            CMDB
          </p>
          <h1 class="text-2xl font-semibold text-slate-900">
            Applications
          </h1>
          <p class="text-sm text-slate-500 mt-1">
            Cartographiez les applications critiques de votre système d'information.
          </p>
        </div>

        <button
            v-if="role === 'Admin'"
            @click="openCreateModal"
            class="bg-slate-900 hover:bg-slate-800 text-white px-4 py-2.5 rounded-xl text-sm font-medium"
        >
          Ajouter une application
        </button>
      </div>

      <div class="grid grid-cols-4 gap-4 mb-6">
        <div class="bg-white border border-slate-100 rounded-2xl p-5">
          <p class="text-sm text-slate-500">Total</p>
          <p class="text-3xl font-bold text-slate-900 mt-1">{{ stats.total }}</p>
        </div>

        <div class="bg-white border border-slate-100 rounded-2xl p-5">
          <p class="text-sm text-slate-500">Actives</p>
          <p class="text-3xl font-bold text-emerald-600 mt-1">{{ stats.active }}</p>
        </div>

        <div class="bg-white border border-slate-100 rounded-2xl p-5">
          <p class="text-sm text-slate-500">Haute criticité</p>
          <p class="text-3xl font-bold text-amber-600 mt-1">{{ stats.high }}</p>
        </div>

        <div class="bg-white border border-slate-100 rounded-2xl p-5">
          <p class="text-sm text-slate-500">Critiques</p>
          <p class="text-3xl font-bold text-red-600 mt-1">{{ stats.critical }}</p>
        </div>
      </div>

      <div class="flex gap-3 mb-5">
        <input
            v-model="searchQuery"
            type="text"
            placeholder="Rechercher application, version, site, owner…"
            class="flex-1 border border-slate-200 rounded-xl px-4 py-2.5 text-sm"
        />

        <select
            v-model="filterCriticality"
            class="border border-slate-200 rounded-xl px-4 py-2.5 text-sm bg-white"
        >
          <option value="all">Toutes criticités</option>
          <option v-for="criticality in criticalityOptions" :key="criticality" :value="criticality">
            {{ criticality }}
          </option>
        </select>
      </div>

      <div v-if="error" class="mb-6 p-4 bg-rose-50 border border-rose-100 text-rose-700 rounded-xl text-sm">
        {{ error }}
      </div>

      <div class="bg-white rounded-2xl border border-slate-100 shadow-sm overflow-hidden">
        <table class="w-full text-left">
          <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="p-4 text-xs uppercase text-slate-400">Application</th>
            <th class="p-4 text-xs uppercase text-slate-400">Version</th>
            <th class="p-4 text-xs uppercase text-slate-400">Criticité</th>
            <th class="p-4 text-xs uppercase text-slate-400">Statut</th>
            <th class="p-4 text-xs uppercase text-slate-400">Site</th>
            <th class="p-4 text-xs uppercase text-slate-400">Owner</th>
            <th class="p-4 text-xs uppercase text-slate-400 text-right">Actions</th>
          </tr>
          </thead>

          <tbody class="divide-y divide-slate-100">
          <tr
              v-for="app in filteredApplications"
              :key="app.ID"
              class="hover:bg-slate-50"
          >
            <td class="p-4">
              <div class="font-medium text-slate-900">{{ app.name }}</div>
              <div class="text-xs text-slate-400 mt-1">{{ app.description || "—" }}</div>
            </td>

            <td class="p-4 text-slate-600">{{ app.version || "—" }}</td>

            <td class="p-4">
                <span
                    class="inline-flex px-2.5 py-1 rounded-full text-xs font-medium border"
                    :class="criticalityConfig[app.criticality]"
                >
                  {{ app.criticality }}
                </span>
            </td>

            <td class="p-4 text-slate-600">{{ app.status || "—" }}</td>
            <td class="p-4 text-slate-600">{{ app.site?.name || "—" }}</td>
            <td class="p-4 text-slate-600">{{ app.owner?.name || "—" }}</td>

            <td class="p-4 text-right">
              <button
                  v-if="role === 'Admin'"
                  @click="openEditModal(app)"
                  class="text-blue-600 hover:underline mr-4"
              >
                Modifier
              </button>

              <button
                  v-if="role === 'Admin'"
                  @click="deleteApplication(app)"
                  class="text-red-600 hover:underline"
              >
                Supprimer
              </button>
            </td>
          </tr>

          <tr v-if="filteredApplications.length === 0">
            <td colspan="7" class="p-10 text-center text-slate-400">
              Aucune application trouvée.
            </td>
          </tr>
          </tbody>
        </table>
      </div>
    </div>
  </DashboardLayout>

  <div
      v-if="showModal"
      class="fixed inset-0 bg-slate-950/40 flex items-center justify-center p-4 z-50"
  >
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-xl">
      <div class="p-6 border-b border-slate-100 flex justify-between">
        <h2 class="text-lg font-semibold">
          {{ isEditMode ? "Modifier l'application" : "Ajouter une application" }}
        </h2>

        <button @click="showModal = false">✕</button>
      </div>

      <form @submit.prevent="saveApplication" class="p-6 space-y-4">
        <input
            v-model="form.name"
            required
            placeholder="Nom de l'application"
            class="input"
        />

        <textarea
            v-model="form.description"
            placeholder="Description"
            class="input"
        ></textarea>

        <div class="grid grid-cols-2 gap-4">
          <input
              v-model="form.version"
              placeholder="Version"
              class="input"
          />

          <select v-model="form.criticality" required class="input">
            <option v-for="criticality in criticalityOptions" :key="criticality" :value="criticality">
              {{ criticality }}
            </option>
          </select>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <select v-model="form.status" class="input">
            <option v-for="status in statusOptions" :key="status" :value="status">
              {{ status }}
            </option>
          </select>

          <select v-model="form.site_id" required class="input">
            <option value="">Sélectionner un site</option>
            <option v-for="site in sites" :key="site.ID" :value="site.ID">
              {{ site.name }}
            </option>
          </select>
        </div>

        <select v-model="form.owner_id" class="input">
          <option value="">Aucun owner</option>
          <option v-for="user in users" :key="user.ID" :value="user.ID">
            {{ user.name }} — {{ user.role }}
          </option>
        </select>

        <div class="flex justify-end gap-3 pt-4">
          <button
              type="button"
              @click="showModal = false"
              class="btn-secondary"
          >
            Annuler
          </button>

          <button
              type="submit"
              :disabled="isSubmitting"
              class="btn-primary"
          >
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