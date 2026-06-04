<script setup>
import { onMounted, ref } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const users = ref([]);
const showModal = ref(false);
const error = ref("");

const form = ref({
  name: "",
  email: "",
  password: "",
  role: "Viewer",
});

const getToken = () => localStorage.getItem("token");

const fetchUsers = async () => {
  try {
    const response = await api.get("/users", {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    users.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Accès refusé ou impossible de charger les utilisateurs";
  }
};

const createUser = async () => {
  try {
    await api.post("/register", form.value);

    showModal.value = false;

    form.value = {
      name: "",
      email: "",
      password: "",
      role: "Viewer",
    };

    fetchUsers();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de créer l'utilisateur";
  }
};

onMounted(fetchUsers);
</script>

<template>
  <DashboardLayout>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">Utilisateurs</h1>

      <button
        @click="showModal = true"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
      >
        Ajouter utilisateur
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
            <th class="p-3 text-left">Email</th>
            <th class="p-3 text-left">Rôle</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="user in users" :key="user.ID" class="border-t">
            <td class="p-3">{{ user.name }}</td>
            <td class="p-3">{{ user.email }}</td>
            <td class="p-3">{{ user.role }}</td>
          </tr>

          <tr v-if="users.length === 0">
            <td colspan="3" class="p-6 text-center text-gray-500">
              Aucun utilisateur trouvé.
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
      <h2 class="text-2xl font-bold mb-6">Ajouter utilisateur</h2>

      <div class="grid gap-4">
        <input v-model="form.name" placeholder="Nom" class="border p-3 rounded-lg" />
        <input v-model="form.email" placeholder="Email" class="border p-3 rounded-lg" />
        <input v-model="form.password" type="password" placeholder="Mot de passe" class="border p-3 rounded-lg" />

        <select v-model="form.role" class="border p-3 rounded-lg">
          <option value="Admin">Admin</option>
          <option value="Technician">Technician</option>
          <option value="Viewer">Viewer</option>
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
          @click="createUser"
          class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
        >
          Créer
        </button>
      </div>
    </div>
  </div>
</template>