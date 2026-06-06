<script setup>
import { onMounted, ref, computed } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const databases = ref([]);
const sites = ref([]);

const showModal = ref(false);
const isEditMode = ref(false);
const editingDatabaseId = ref(null);

const error = ref("");
const isSubmitting = ref(false);
const searchQuery = ref("");
const filterEngine = ref("all");

const form = ref({
  name: "",
  engine: "PostgreSQL",
  version: "",
  environment: "Production",
  site_id: "",
});

const getToken = () => localStorage.getItem("token");
const role = computed(() => localStorage.getItem("role") || "");

const engineOptions = ["PostgreSQL", "MySQL", "SQL Server", "Oracle", "MongoDB", "MariaDB", "Redis", "Autre"];
const environmentOptions = ["Production", "Preproduction", "Staging", "Development", "Test"];

const engineBadgeClass = (engine) => {
  switch (engine) {
    case "PostgreSQL":
      return "bg-blue-50 text-blue-700 border-blue-100";
    case "MySQL":
    case "MariaDB":
      return "bg-orange-50 text-orange-700 border-orange-100";
    case "SQL Server":
      return "bg-red-50 text-red-700 border-red-100";
    case "MongoDB":
      return "bg-green-50 text-green-700 border-green-100";
    case "Redis":
      return "bg-rose-50 text-rose-700 border-rose-100";
    default:
      return "bg-slate-100 text-slate-700 border-slate-200";
  }
};

const environmentBadgeClass = (environment) => {
  switch (environment) {
    case "Production":
      return "bg-red-50 text-red-700 border-red-100";
    case "Preproduction":
    case "Staging":
      return "bg-amber-50 text-amber-700 border-amber-100";
    case "Development":
    case "Test":
      return "bg-slate-100 text-slate-700 border-slate-200";
    default:
      return "bg-slate-100 text-slate-700 border-slate-200";
  }
};

const filteredDatabases = computed(() => {
  let list = databases.value;

  if (filterEngine.value !== "all") {
    list = list.filter((db) => db.engine === filterEngine.value);
  }

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();

    list = list.filter(
        (db) =>
            db.name?.toLowerCase().includes(q) ||
            db.engine?.toLowerCase().includes(q) ||
            db.version?.toLowerCase().includes(q) ||
            db.environment?.toLowerCase().includes(q) ||
            db.site?.name?.toLowerCase().includes(q)
    );
  }

  return list;
});

const stats = computed(() => ({
  total: databases.value.length,
  production: databases.value.filter((db) => db.environment === "Production").length,
  postgres: databases.value.filter((db) => db.engine === "PostgreSQL").length,
  mysql: databases.value.filter((db) => db.engine === "MySQL").length,
}));

const fetchDatabases = async () => {
  try {
    const response = await api.get("/cmdb/databases", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });

    databases.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger les bases de données";
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

const resetForm = () => {
  form.value = {
    name: "",
    engine: "PostgreSQL",
    version: "",
    environment: "Production",
    site_id: "",
  };

  isEditMode.value = false;
  editingDatabaseId.value = null;
  error.value = "";
};

const openCreateModal = () => {
  resetForm();
  showModal.value = true;
};

const openEditModal = (db) => {
  isEditMode.value = true;
  editingDatabaseId.value = db.ID;

  form.value = {
    name: db.name,
    engine: db.engine,
    version: db.version,
    environment: db.environment,
    site_id: db.site_id,
  };

  showModal.value = true;
};

const saveDatabase = async () => {
  try {
    isSubmitting.value = true;
    error.value = "";

    const payload = {
      ...form.value,
      site_id: form.value.site_id ? Number(form.value.site_id) : null,
    };

    if (isEditMode.value) {
      await api.put(`/cmdb/databases/${editingDatabaseId.value}`, payload, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    } else {
      await api.post("/cmdb/databases", payload, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    }

    showModal.value = false;
    resetForm();
    await fetchDatabases();
  } catch (err) {
    console.error(err);
    error.value = "Impossible d'enregistrer la base de données";
  } finally {
    isSubmitting.value = false;
  }
};

const deleteDatabase = async (db) => {
  if (!confirm(`Supprimer la base "${db.name}" ?`)) return;

  try {
    await api.delete(`/cmdb/databases/${db.ID}`, {
      headers: { Authorization: `Bearer ${getToken()}` },
    });

    await fetchDatabases();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de supprimer la base de données";
  }
};

onMounted(() => {
  fetchDatabases();
  fetchSites();
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
            Bases de données
          </h1>
          <p class="text-sm text-slate-500 mt-1">
            Inventoriez les bases critiques et leur localisation.
          </p>
        </div>

        <button
            v-if="role === 'Admin'"
            @click="openCreateModal"
            class="bg-slate-900 hover:bg-slate-800 text-white px-4 py-2.5 rounded-xl text-sm font-medium"
        >
          Ajouter une base
        </button>
      </div>

      <div class="grid grid-cols-4 gap-4 mb-6">
        <div class="bg-white border border-slate-100 rounded-2xl p-5">
          <p class="text-sm text-slate-500">Total</p>
          <p class="text-3xl font-bold text-slate-900 mt-1">{{ stats.total }}</p>
        </div>

        <div class="bg-white border border-slate-100 rounded-2xl p-5">
          <p class="text-sm text-slate-500">Production</p>
          <p class="text-3xl font-bold text-red-600 mt-1">{{ stats.production }}</p>
        </div>

        <div class="bg-white border border-slate-100 rounded-2xl p-5">
          <p class="text-sm text-slate-500">PostgreSQL</p>
          <p class="text-3xl font-bold text-blue-600 mt-1">{{ stats.postgres }}</p>
        </div>

        <div class="bg-white border border-slate-100 rounded-2xl p-5">
          <p class="text-sm text-slate-500">MySQL</p>
          <p class="text-3xl font-bold text-orange-600 mt-1">{{ stats.mysql }}</p>
        </div>
      </div>

      <div class="flex gap-3 mb-5">
        <input
            v-model="searchQuery"
            type="text"
            placeholder="Rechercher base, moteur, version, environnement, site…"
            class="flex-1 border border-slate-200 rounded-xl px-4 py-2.5 text-sm"
        />

        <select
            v-model="filterEngine"
            class="border border-slate-200 rounded-xl px-4 py-2.5 text-sm bg-white"
        >
          <option value="all">Tous les moteurs</option>
          <option v-for="engine in engineOptions" :key="engine" :value="engine">
            {{ engine }}
          </option>
        </select>
      </div>

      <div
          v-if="error"
          class="mb-6 p-4 bg-rose-50 border border-rose-100 text-rose-700 rounded-xl text-sm"
      >
        {{ error }}
      </div>

      <div class="bg-white rounded-2xl border border-slate-100 shadow-sm overflow-hidden">
        <table class="w-full text-left">
          <thead class="bg-slate-50 border-b border-slate-100">
          <tr>
            <th class="p-4 text-xs uppercase text-slate-400">Nom</th>
            <th class="p-4 text-xs uppercase text-slate-400">Moteur</th>
            <th class="p-4 text-xs uppercase text-slate-400">Version</th>
            <th class="p-4 text-xs uppercase text-slate-400">Environnement</th>
            <th class="p-4 text-xs uppercase text-slate-400">Site</th>
            <th class="p-4 text-xs uppercase text-slate-400 text-right">Actions</th>
          </tr>
          </thead>

          <tbody class="divide-y divide-slate-100">
          <tr
              v-for="db in filteredDatabases"
              :key="db.ID"
              class="hover:bg-slate-50"
          >
            <td class="p-4 font-medium text-slate-900">
              {{ db.name }}
            </td>

            <td class="p-4">
                <span
                    class="inline-flex px-2.5 py-1 rounded-full text-xs font-medium border"
                    :class="engineBadgeClass(db.engine)"
                >
                  {{ db.engine }}
                </span>
            </td>

            <td class="p-4 text-slate-600">
              {{ db.version || "—" }}
            </td>

            <td class="p-4">
                <span
                    class="inline-flex px-2.5 py-1 rounded-full text-xs font-medium border"
                    :class="environmentBadgeClass(db.environment)"
                >
                  {{ db.environment || "—" }}
                </span>
            </td>

            <td class="p-4 text-slate-600">
              {{ db.site?.name || "—" }}
            </td>

            <td class="p-4 text-right">
              <button
                  v-if="role === 'Admin'"
                  @click="openEditModal(db)"
                  class="text-blue-600 hover:underline mr-4"
              >
                Modifier
              </button>

              <button
                  v-if="role === 'Admin'"
                  @click="deleteDatabase(db)"
                  class="text-red-600 hover:underline"
              >
                Supprimer
              </button>
            </td>
          </tr>

          <tr v-if="filteredDatabases.length === 0">
            <td colspan="6" class="p-10 text-center text-slate-400">
              Aucune base de données trouvée.
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
          {{ isEditMode ? "Modifier la base" : "Ajouter une base" }}
        </h2>

        <button @click="showModal = false">✕</button>
      </div>

      <form @submit.prevent="saveDatabase" class="p-6 space-y-4">
        <input
            v-model="form.name"
            required
            placeholder="Nom de la base"
            class="input"
        />

        <div class="grid grid-cols-2 gap-4">
          <select v-model="form.engine" required class="input">
            <option v-for="engine in engineOptions" :key="engine" :value="engine">
              {{ engine }}
            </option>
          </select>

          <input
              v-model="form.version"
              placeholder="Version"
              class="input"
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <select v-model="form.environment" class="input">
            <option v-for="environment in environmentOptions" :key="environment" :value="environment">
              {{ environment }}
            </option>
          </select>

          <select v-model="form.site_id" class="input">
            <option value="">Sélectionner un site</option>
            <option v-for="site in sites" :key="site.ID" :value="site.ID">
              {{ site.name }}
            </option>
          </select>
        </div>

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