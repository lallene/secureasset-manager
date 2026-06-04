<script setup>
import { onMounted, ref } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const assets = ref([]);
const showModal = ref(false);

const form = ref({
  name: "",
  type: "",
  serial: "",
  ip_address: "",
  site: "",
  status: "",
  assigned_to: "",
});

const fetchAssets = async () => {
  try {
    const token = localStorage.getItem("token");

    const response = await api.get("/assets", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    assets.value = response.data;
  } catch (error) {
    console.error(error);
  }
};

const createAsset = async () => {
  try {
    const token = localStorage.getItem("token");

    await api.post("/assets", form.value, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    showModal.value = false;

    form.value = {
      name: "",
      type: "",
      serial: "",
      ip_address: "",
      site: "",
      status: "",
      assigned_to: "",
    };

    fetchAssets();
  } catch (error) {
    console.error(error);
  }
};

onMounted(fetchAssets);
</script>

<template>
  <DashboardLayout>
    <h1 class="text-2xl font-bold mb-6">Assets</h1>

    <div class="mb-6">
      <button
        @click="showModal = true"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
      >
        Ajouter Asset
      </button>
    </div>

    <div class="bg-white rounded-xl shadow overflow-hidden">
      <table class="w-full">
        <thead class="bg-slate-100">
          <tr>
            <th class="p-3 text-left">Nom</th>
            <th class="p-3 text-left">Type</th>
            <th class="p-3 text-left">IP</th>
            <th class="p-3 text-left">Site</th>
            <th class="p-3 text-left">Statut</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="asset in assets" :key="asset.ID" class="border-t">
            <td class="p-3">{{ asset.name }}</td>
            <td class="p-3">{{ asset.type }}</td>
            <td class="p-3">{{ asset.ip_address }}</td>
            <td class="p-3">{{ asset.site }}</td>
            <td class="p-3">{{ asset.status }}</td>
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
        Ajouter Asset
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
          @click="createAsset"
          class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
        >
          Créer
        </button>
      </div>
    </div>
  </div>
</template>