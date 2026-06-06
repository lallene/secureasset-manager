<script setup>
import { onMounted, ref } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const users = ref([]);
const showModal = ref(false);
const error = ref("");
const isSubmitting = ref(false);

const sites = ref([]);
const services = ref([]);
const isEditMode = ref(false);
const editingUserId = ref(null);

const form = ref({
  name: "",
  email: "",
  password: "",
  role: "Requester",
  site_id: "",
  service_id: "",
});

const fetchSites = async () => {
  const response = await api.get("/sites", {
    headers: { Authorization: `Bearer ${getToken()}` },
  });
  sites.value = response.data;
};

const fetchServices = async () => {
  const response = await api.get("/services", {
    headers: { Authorization: `Bearer ${getToken()}` },
  });
  services.value = response.data;
};

const getToken = () => localStorage.getItem("token");

const fetchUsers = async () => {
  try {
    const response = await api.get("/users", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    users.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Accès refusé ou impossible de charger les utilisateurs";
  }
};


const saveUser = async () => {
  try {
    isSubmitting.value = true;
    error.value = "";

    const payload = {
      ...form.value,
      site_id: Number(form.value.site_id),
      service_id: form.value.service_id ? Number(form.value.service_id) : null,
    };

    if (isEditMode.value) {
      await api.put(`/users/${editingUserId.value}`, payload, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    } else {
      await api.post("/users", payload, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    }

    showModal.value = false;
    resetForm();
    await fetchUsers();
  } catch (err) {
    console.error(err);
    error.value = "Impossible d'enregistrer l'utilisateur";
  } finally {
    isSubmitting.value = false;
  }
};

const resetForm = () => {
  form.value = {
    name: "",
    email: "",
    password: "",
    role: "Requester",
    site_id: "",
    service_id: "",
  };

  isEditMode.value = false;
  editingUserId.value = null;
};

const openCreateModal = () => {
  resetForm();
  showModal.value = true;
};

const openEditModal = (user) => {
  isEditMode.value = true;
  editingUserId.value = user.ID;

  form.value = {
    name: user.name,
    email: user.email,
    password: "",
    role: user.role,
    site_id: user.site_id,
    service_id: user.service_id || "",
  };

  showModal.value = true;
};

const deleteUser = async (user) => {
  if (!confirm(`Supprimer l'utilisateur ${user.name} ?`)) return;

  await api.delete(`/users/${user.ID}`, {
    headers: { Authorization: `Bearer ${getToken()}` },
  });

  fetchUsers();
};



// Utilitaire de style pour afficher des badges de rôle élégants
const getRoleBadgeClass = (role) => {
  switch (role?.toLowerCase()) {
    case "admin":
      return "bg-purple-50 text-purple-700 border-purple-100";

    case "agent":
      return "bg-blue-50 text-blue-700 border-blue-100";

    case "requester":
      return "bg-green-50 text-green-700 border-green-100";

    default:
      return "bg-slate-100 text-slate-700 border-slate-200";
  }
};

onMounted(() => {
  fetchUsers();
  fetchSites();
  fetchServices();
});
</script>

<template>
  <DashboardLayout class="bg-slate-50/50 min-h-screen text-slate-900">
    <!-- Entête de page -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight text-slate-900">Utilisateurs</h1>
        <p class="text-sm text-slate-500 mt-1">Gérez les accès et les permissions de vos collaborateurs.</p>
      </div>

      <button
          @click="openCreateModal"
        class="inline-flex items-center gap-2 bg-slate-900 hover:bg-slate-800 text-white font-medium text-sm px-4 py-2.5 rounded-xl transition-all duration-150 shadow-sm"
      >
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
        </svg>
        Ajouter un utilisateur
      </button>
    </div>

    <!-- Alerte d'erreur -->
    <div v-if="error" class="mb-6 p-4 bg-rose-50 border border-rose-100 text-rose-700 rounded-xl text-sm flex items-center gap-2">
      <span class="w-1.5 h-1.5 rounded-full bg-rose-500"></span>
      {{ error }}
    </div>

    <!-- Table d'utilisateurs Premium -->
    <div class="bg-white rounded-2xl border border-slate-100 shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full border-collapse text-left">
          <thead>
            <tr class="border-b border-slate-100 bg-slate-50/70">
              <th class="p-4 text-xs font-semibold uppercase tracking-wider text-slate-400">Nom</th>
              <th class="p-4 text-xs font-semibold uppercase tracking-wider text-slate-400">Adresse Email</th>
              <th class="p-4 text-xs font-semibold uppercase tracking-wider text-slate-400">Sites</th>
              <th class="p-4 text-xs font-semibold uppercase tracking-wider text-slate-400">Services</th>
              <th class="p-4 text-xs font-semibold uppercase tracking-wider text-slate-400">Rôle</th>
            </tr>
          </thead>

          <tbody class="divide-y divide-slate-100 text-sm text-slate-700">
            <tr v-for="user in users" :key="user.ID" class="hover:bg-slate-50/40 transition-colors group">
              <!-- Colonne Nom avec un léger accent gras -->
              <td class="p-4 font-medium text-slate-900">{{ user.name }}</td>
              <td class="p-4 text-slate-500">{{ user.email }}</td>
              <td class="p-4 text-slate-500">{{ user.site?.name || "-" }}</td>
              <td class="p-4 text-slate-500">{{ user.service?.name || "-" }}</td>
              <td class="p-4">
                <span :class="getRoleBadgeClass(user.role)" class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium border">
                  {{ user.role }}
                </span>
              </td>
              <td class="flex gap-2">
                <button
                    @click="openEditModal(user)"
                    class="text-blue-600 hover:text-blue-800"
                >
                  Modifier
                </button>

                <button
                    @click="deleteUser(user)"
                    class="text-red-600 hover:text-red-800"
                >
                  Supprimer
                </button>
              </td>
            </tr>

            <!-- État table vide -->
            <tr v-if="users.length === 0">
              <td colspan="3" class="p-12 text-center text-slate-400">
                <svg class="w-8 h-8 text-slate-300 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
                </svg>
                Aucun utilisateur trouvé.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </DashboardLayout>

  <!-- Modale Royale (Overlay fixe avec flou d'arrière-plan) -->
  <div
    v-if="showModal"
    class="fixed inset-0 bg-slate-950/40 backdrop-blur-sm flex items-center justify-center p-4 z-50 animate-fadeIn"
  >
    <div class="bg-white rounded-2xl border border-slate-100 shadow-2xl w-full max-w-lg overflow-hidden transform transition-all">
      <!-- Header de la modale -->
      <div class="px-6 py-4 border-b border-slate-100 flex items-center justify-between">
        <h2 class="text-lg font-semibold text-slate-900">Ajouter un collaborateur</h2>
        <button @click="showModal = false" class="text-slate-400 hover:text-slate-600 rounded-lg p-1 hover:bg-slate-50 transition-colors">
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Formulaire natif -->
      <form @submit.prevent="saveUser">
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1.5">Nom complet</label>
            <input v-model="form.name" required placeholder="Ex: Jean Dupont" class="w-full border border-slate-200 rounded-xl px-3 py-2.5 text-sm placeholder-slate-400 focus:outline-none focus:border-blue-600 focus:ring-4 focus:ring-blue-50 transition-all bg-slate-50/40 focus:bg-white" />
          </div>

          <div>
            <label class="block text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1.5">Adresse Email</label>
            <input v-model="form.email" type="email" required placeholder="j.dupont@entreprise.com" class="w-full border border-slate-200 rounded-xl px-3 py-2.5 text-sm placeholder-slate-400 focus:outline-none focus:border-blue-600 focus:ring-4 focus:ring-blue-50 transition-all bg-slate-50/40 focus:bg-white" />
          </div>

          <div>
            <label class="block text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1.5">Mot de passe temporaire</label>
            <input
                v-model="form.password"
                type="password"
                :required="!isEditMode"
                placeholder="Laisser vide pour conserver le mot de passe"
                class="w-full border border-slate-200 rounded-xl px-3 py-2.5 text-sm"
            />
          </div>
          <select v-model="form.site_id" required>
            <option value="">Sélectionner un site</option>
            <option v-for="site in sites" :key="site.ID" :value="site.ID">
              {{ site.name }}
            </option>
          </select>
          <select v-model="form.service_id">
            <option value="">Aucun service</option>
            <option v-for="service in services" :key="service.ID" :value="service.ID">
              {{ service.name }}
            </option>
          </select>

          <div>
            <label class="block text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1.5">Rôle & Permissions</label>
            <div class="relative">
              <select v-model="form.role" class="w-full border border-slate-200 rounded-xl px-3 py-2.5 text-sm bg-slate-50/40 focus:bg-white focus:outline-none focus:border-blue-600 focus:ring-4 focus:ring-blue-50 transition-all appearance-none cursor-pointer text-slate-700">
                <option value="Admin">Admin</option>
                <option value="Agent">Agent</option>
                <option value="Requester">Requester</option>
              </select>
              <!-- Petite flèche personnalisée pour remplacer le sélecteur natif moche -->
              <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none text-slate-400">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                </svg>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer actions de la modale -->
        <div class="px-6 py-4 bg-slate-50/80 border-t border-slate-100 flex justify-end gap-3">
          <button
            type="button"
            @click="showModal = false"
            class="px-4 py-2 border border-slate-200 hover:bg-slate-100 text-slate-600 font-medium text-sm rounded-xl transition-colors"
          >
            Annuler
          </button>

          <button
            type="submit"
            :disabled="isSubmitting"
            class="bg-slate-900 hover:bg-slate-800 text-white px-4 py-2 font-medium text-sm rounded-xl transition-all shadow-sm flex items-center gap-2 disabled:opacity-50"
          >
            <span v-if="isSubmitting" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <span>Créer le compte</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
.animate-fadeIn {
  animation: fadeIn 0.15s ease-out forwards;
}
</style>