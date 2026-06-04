<script setup>
import { onMounted, ref } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const assets = ref([]);
const showModal = ref(false);
const isEditMode = ref(false);
const editingAssetId = ref(null);
const error = ref("");

const form = ref({
  name: "",
  type: "",
  serial: "",
  ip_address: "",
  site: "",
  status: "",
  assigned_to: "",
});

const getToken = () => localStorage.getItem("token");

const resetForm = () => {
  form.value = {
    name: "",
    type: "",
    serial: "",
    ip_address: "",
    site: "",
    status: "",
    assigned_to: "",
  };

  isEditMode.value = false;
  editingAssetId.value = null;
  error.value = "";
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
    error.value = "Impossible de charger les assets";
  }
};

const openCreateModal = () => {
  resetForm();
  showModal.value = true;
};

const openEditModal = (asset) => {
  isEditMode.value = true;
  editingAssetId.value = asset.ID;

  form.value = {
    name: asset.name,
    type: asset.type,
    serial: asset.serial,
    ip_address: asset.ip_address,
    site: asset.site,
    status: asset.status,
    assigned_to: asset.assigned_to,
  };

  showModal.value = true;
};

const saveAsset = async () => {
  try {
    if (isEditMode.value) {
      await api.put(`/assets/${editingAssetId.value}`, form.value, {
        headers: {
          Authorization: `Bearer ${getToken()}`,
        },
      });
    } else {
      await api.post("/assets", form.value, {
        headers: {
          Authorization: `Bearer ${getToken()}`,
        },
      });
    }

    showModal.value = false;
    resetForm();
    fetchAssets();
  } catch (err) {
    console.error(err);
    error.value = "Erreur lors de l'enregistrement";
  }
};

const deleteAsset = async (asset) => {
  const confirmed = confirm(`Supprimer l'asset ${asset.name} ?`);

  if (!confirmed) {
    return;
  }

  try {
    await api.delete(`/assets/${asset.ID}`, {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    fetchAssets();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de supprimer cet asset";
  }
};

onMounted(fetchAssets);
</script>

<template>
  <DashboardLayout>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">Assets</h1>

      <button
        @click="openCreateModal"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
      >
        Ajouter Asset
      </button>
    </div>

    <p v-if="error" class="mb-4 text-red-600">
      {{ error }}
    </p>

    <div class="bg-white rounded-xl shadow overflow-hidden">
      <table class="w-full">
        <thead class="bg-slate-100">
          <tr>
            <th class="p-3 text-left">Nom</th>
            <th class="p-3 text-left">Type</th>
            <th class="p-3 text-left">IP</th>
            <th class="p-3 text-left">Site</th>
            <th class="p-3 text-left">Statut</th>
            <th class="p-3 text-right">Actions</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="asset in assets" :key="asset.ID" class="border-t">
            <td class="p-3">{{ asset.name }}</td>
            <td class="p-3">{{ asset.type }}</td>
            <td class="p-3">{{ asset.ip_address }}</td>
            <td class="p-3">{{ asset.site }}</td>
            <td class="p-3">{{ asset.status }}</td>
            <td class="p-3 text-right">
              <button
                @click="openEditModal(asset)"
                class="text-blue-600 hover:underline mr-4"
              >
                Modifier
              </button>

              <button
                v-if="role === 'Admin'"
                @click="deleteAsset(asset)"
              >
                Supprimer
              </button>
            </td>
          </tr>

          <tr v-if="assets.length === 0">
            <td colspan="6" class="p-6 text-center text-gray-500">
              Aucun asset trouvé.
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
        {{ isEditMode ? "Modifier Asset" : "Ajouter Asset" }}
      </h2>

      <div class="grid gap-4">
        <input v-model="form.name" placeholder="Nom" class="border p-3 rounded-lg" />
        <input v-model="form.type" placeholder="Type" class="border p-3 rounded-lg" />
        <input v-model="form.serial" placeholder="Serial" class="border p-3 rounded-lg" />
        <input v-model="form.ip_address" placeholder="Adresse IP" class="border p-3 rounded-lg" />
        <input v-model="form.site" placeholder="Site" class="border p-3 rounded-lg" />
        <input v-model="form.status" placeholder="Status" class="border p-3 rounded-lg" />
        <input v-model="form.assigned_to" placeholder="Assigné à" class="border p-3 rounded-lg" />
      </div>

      <div class="flex justify-end gap-4 mt-6">
        <button
          @click="showModal = false"
          class="px-4 py-2 border rounded-lg"
        >
          Annuler
        </button>

        <button
          @click="saveAsset"
          class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
        >
          {{ isEditMode ? "Enregistrer" : "Créer" }}
        </button>
      </div>
    </div>
  </div>
</template>