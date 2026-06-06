<script setup>
import { onMounted, ref, computed } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const relations = ref([]);
const applications = ref([]);
const databases = ref([]);
const assets = ref([]);

const showModal = ref(false);
const isEditMode = ref(false);
const editingRelationId = ref(null);

const error = ref("");
const isSubmitting = ref(false);
const searchQuery = ref("");

const form = ref({
  source_type: "Application",
  source_id: "",
  relation_type: "hosted_on",
  target_type: "Asset",
  target_id: "",
});

const getToken = () => localStorage.getItem("token");
const role = computed(() => localStorage.getItem("role") || "");

const typeOptions = ["Application", "Asset", "Database"];
const relationOptions = ["hosted_on", "uses_database", "depends_on", "connects_to"];

const headers = () => ({
  headers: { Authorization: `Bearer ${getToken()}` },
});

const fetchRelations = async () => {
  try {
    const response = await api.get("/cmdb/relations", headers());
    relations.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger les relations CMDB";
  }
};

const fetchApplications = async () => {
  try {
    const response = await api.get("/cmdb/applications", headers());
    applications.value = response.data;
  } catch (err) {
    console.error(err);
  }
};

const fetchDatabases = async () => {
  try {
    const response = await api.get("/cmdb/databases", headers());
    databases.value = response.data;
  } catch (err) {
    console.error(err);
  }
};

const fetchAssets = async () => {
  try {
    const response = await api.get("/assets", headers());
    assets.value = response.data;
  } catch (err) {
    console.error(err);
  }
};

const getItemsByType = (type) => {
  switch (type) {
    case "Application":
      return applications.value;
    case "Asset":
      return assets.value;
    case "Database":
      return databases.value;
    default:
      return [];
  }
};

const getItemLabel = (type, id) => {
  const item = getItemsByType(type).find((el) => Number(el.ID) === Number(id));

  if (!item) {
    return `${type} #${id}`;
  }

  if (type === "Asset") {
    return item.name || `Asset #${id}`;
  }

  if (type === "Database") {
    return item.name || `Database #${id}`;
  }

  return item.name || `${type} #${id}`;
};

const resetForm = () => {
  form.value = {
    source_type: "Application",
    source_id: "",
    relation_type: "hosted_on",
    target_type: "Asset",
    target_id: "",
  };

  isEditMode.value = false;
  editingRelationId.value = null;
  error.value = "";
};

const openCreateModal = () => {
  resetForm();
  showModal.value = true;
};

const openEditModal = (relation) => {
  isEditMode.value = true;
  editingRelationId.value = relation.ID;

  form.value = {
    source_type: relation.source_type,
    source_id: relation.source_id,
    relation_type: relation.relation_type,
    target_type: relation.target_type,
    target_id: relation.target_id,
  };

  showModal.value = true;
};

const saveRelation = async () => {
  try {
    isSubmitting.value = true;
    error.value = "";

    const payload = {
      ...form.value,
      source_id: Number(form.value.source_id),
      target_id: Number(form.value.target_id),
    };

    if (isEditMode.value) {
      await api.put(`/cmdb/relations/${editingRelationId.value}`, payload, headers());
    } else {
      await api.post("/cmdb/relations", payload, headers());
    }

    showModal.value = false;
    resetForm();
    await fetchRelations();
  } catch (err) {
    console.error(err);
    error.value = "Impossible d'enregistrer la relation";
  } finally {
    isSubmitting.value = false;
  }
};

const deleteRelation = async (relation) => {
  if (!confirm("Supprimer cette relation CMDB ?")) return;

  try {
    await api.delete(`/cmdb/relations/${relation.ID}`, headers());
    await fetchRelations();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de supprimer la relation";
  }
};

const filteredRelations = computed(() => {
  if (!searchQuery.value.trim()) {
    return relations.value;
  }

  const q = searchQuery.value.toLowerCase();

  return relations.value.filter((rel) => {
    const source = getItemLabel(rel.source_type, rel.source_id).toLowerCase();
    const target = getItemLabel(rel.target_type, rel.target_id).toLowerCase();

    return (
        source.includes(q) ||
        target.includes(q) ||
        rel.relation_type?.toLowerCase().includes(q) ||
        rel.source_type?.toLowerCase().includes(q) ||
        rel.target_type?.toLowerCase().includes(q)
    );
  });
});

const stats = computed(() => ({
  total: relations.value.length,
  hosted: relations.value.filter((r) => r.relation_type === "hosted_on").length,
  database: relations.value.filter((r) => r.relation_type === "uses_database").length,
  depends: relations.value.filter((r) => r.relation_type === "depends_on").length,
}));

onMounted(() => {
  fetchApplications();
  fetchDatabases();
  fetchAssets();
  fetchRelations();
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
            Relations
          </h1>
          <p class="text-sm text-slate-500 mt-1">
            Cartographiez les dépendances entre applications, assets et bases de données.
          </p>
        </div>

        <button
            v-if="role === 'Admin'"
            @click="openCreateModal"
            class="bg-slate-900 hover:bg-slate-800 text-white px-4 py-2.5 rounded-xl text-sm font-medium"
        >
          Ajouter une relation
        </button>
      </div>

      <div class="grid grid-cols-4 gap-4 mb-6">
        <div class="bg-white border border-slate-100 rounded-2xl p-5">
          <p class="text-sm text-slate-500">Relations</p>
          <p class="text-3xl font-bold text-slate-900 mt-1">{{ stats.total }}</p>
        </div>

        <div class="bg-white border border-slate-100 rounded-2xl p-5">
          <p class="text-sm text-slate-500">Hébergement</p>
          <p class="text-3xl font-bold text-indigo-600 mt-1">{{ stats.hosted }}</p>
        </div>

        <div class="bg-white border border-slate-100 rounded-2xl p-5">
          <p class="text-sm text-slate-500">Bases liées</p>
          <p class="text-3xl font-bold text-blue-600 mt-1">{{ stats.database }}</p>
        </div>

        <div class="bg-white border border-slate-100 rounded-2xl p-5">
          <p class="text-sm text-slate-500">Dépendances</p>
          <p class="text-3xl font-bold text-amber-600 mt-1">{{ stats.depends }}</p>
        </div>
      </div>

      <div class="mb-5">
        <input
            v-model="searchQuery"
            type="text"
            placeholder="Rechercher source, cible ou type de relation…"
            class="w-full border border-slate-200 rounded-xl px-4 py-2.5 text-sm"
        />
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
            <th class="p-4 text-xs uppercase text-slate-400">Source</th>
            <th class="p-4 text-xs uppercase text-slate-400">Relation</th>
            <th class="p-4 text-xs uppercase text-slate-400">Cible</th>
            <th class="p-4 text-xs uppercase text-slate-400 text-right">Actions</th>
          </tr>
          </thead>

          <tbody class="divide-y divide-slate-100">
          <tr
              v-for="relation in filteredRelations"
              :key="relation.ID"
              class="hover:bg-slate-50"
          >
            <td class="p-4">
              <div class="font-medium text-slate-900">
                {{ getItemLabel(relation.source_type, relation.source_id) }}
              </div>
              <div class="text-xs text-slate-400 mt-1">
                {{ relation.source_type }}
              </div>
            </td>

            <td class="p-4">
                <span class="inline-flex px-2.5 py-1 rounded-full text-xs font-medium border bg-indigo-50 text-indigo-700 border-indigo-100">
                  {{ relation.relation_type }}
                </span>
            </td>

            <td class="p-4">
              <div class="font-medium text-slate-900">
                {{ getItemLabel(relation.target_type, relation.target_id) }}
              </div>
              <div class="text-xs text-slate-400 mt-1">
                {{ relation.target_type }}
              </div>
            </td>

            <td class="p-4 text-right">
              <button
                  v-if="role === 'Admin'"
                  @click="openEditModal(relation)"
                  class="text-blue-600 hover:underline mr-4"
              >
                Modifier
              </button>

              <button
                  v-if="role === 'Admin'"
                  @click="deleteRelation(relation)"
                  class="text-red-600 hover:underline"
              >
                Supprimer
              </button>
            </td>
          </tr>

          <tr v-if="filteredRelations.length === 0">
            <td colspan="4" class="p-10 text-center text-slate-400">
              Aucune relation CMDB trouvée.
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
          {{ isEditMode ? "Modifier la relation" : "Ajouter une relation" }}
        </h2>

        <button @click="showModal = false">✕</button>
      </div>

      <form @submit.prevent="saveRelation" class="p-6 space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <select v-model="form.source_type" class="input" @change="form.source_id = ''">
            <option v-for="type in typeOptions" :key="type" :value="type">
              {{ type }}
            </option>
          </select>

          <select v-model="form.source_id" required class="input">
            <option value="">Sélectionner la source</option>
            <option
                v-for="item in getItemsByType(form.source_type)"
                :key="item.ID"
                :value="item.ID"
            >
              {{ item.name }}
            </option>
          </select>
        </div>

        <select v-model="form.relation_type" required class="input">
          <option v-for="relation in relationOptions" :key="relation" :value="relation">
            {{ relation }}
          </option>
        </select>

        <div class="grid grid-cols-2 gap-4">
          <select v-model="form.target_type" class="input" @change="form.target_id = ''">
            <option v-for="type in typeOptions" :key="type" :value="type">
              {{ type }}
            </option>
          </select>

          <select v-model="form.target_id" required class="input">
            <option value="">Sélectionner la cible</option>
            <option
                v-for="item in getItemsByType(form.target_type)"
                :key="item.ID"
                :value="item.ID"
            >
              {{ item.name }}
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