<script setup>
import { onMounted, ref, computed } from "vue";
import { useRouter } from "vue-router";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const notifications = ref([]);
const error = ref("");

const getToken = () => localStorage.getItem("token");

const unreadCount = computed(() => {
  return notifications.value.filter((n) => !n.is_read).length;
});

const router = useRouter();

const goToIncident = async (notification) => {
  if (!notification.incident_id) return;
  if (!notification.is_read) {
    await markAsRead(notification);
  }
  router.push(`/incidents?incident=${notification.incident_id}`);
};

const fetchNotifications = async () => {
  try {
    const response = await api.get("/notifications", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    notifications.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger les notifications";
  }
};

const markAsRead = async (notification) => {
  try {
    await api.put(`/notifications/${notification.ID}/read`, {}, {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    notification.is_read = true;
  } catch (err) {
    console.error(err);
  }
};

const markAllAsRead = async () => {
  // Optionnel mais très pro : boucle ou endpoint pour tout marquer d'un coup
  const unreadList = notifications.value.filter(n => !n.is_read);
  for (const n of unreadList) {
    await markAsRead(n);
  }
};

onMounted(fetchNotifications);
</script>

<template>
  <DashboardLayout class="bg-slate-50/50 min-h-screen text-slate-900">
    <!-- Header de la page -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight text-slate-900">Notifications</h1>
        <p class="text-sm text-slate-500 mt-1">Restez informé des mises à jour sur vos incidents.</p>
      </div>

      <div class="flex items-center gap-3">
        <button
          v-if="unreadCount > 0"
          @click="markAllAsRead"
          class="text-xs font-medium text-slate-500 hover:text-slate-800 transition-colors px-3 py-2 rounded-lg hover:bg-slate-100"
        >
          Tout marquer comme lu
        </button>
        <span
          :class="unreadCount > 0 ? 'bg-blue-50 text-blue-700 border-blue-100' : 'bg-slate-100 text-slate-600 border-slate-200'"
          class="px-3 py-1 rounded-full text-xs font-semibold border"
        >
          {{ unreadCount }} non lues
        </span>
      </div>
    </div>

    <!-- Alert d'erreur épurée -->
    <div v-if="error" class="mb-6 p-4 bg-rose-50 border border-rose-100 text-rose-700 rounded-xl text-sm flex items-center gap-2">
      <span class="w-1.5 h-1.5 rounded-full bg-rose-500"></span>
      {{ error }}
    </div>

    <!-- Liste des notifications -->
    <div class="space-y-3 max-w-4xl">
      <div
        v-for="notification in notifications"
        :key="notification.ID"
        class="bg-white rounded-xl border transition-all duration-200 p-5 group relative"
        :class="notification.is_read
          ? 'border-slate-100 bg-white/60'
          : 'border-slate-200 shadow-sm bg-gradient-to-r from-blue-50/30 to-transparent'"
      >
        <div class="flex items-start gap-4">
          <!-- Indicateur de lecture (Pastille bleue ou grise fine) -->
          <div class="mt-1.5 flex-shrink-0">
            <span
              v-if="!notification.is_read"
              class="block w-2.5 h-2.5 rounded-full bg-blue-600 ring-4 ring-blue-50"
            ></span>
            <span
              v-else
              class="block w-2.5 h-2.5 rounded-full bg-slate-200"
            ></span>
          </div>

          <!-- Contenu -->
          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-4">
              <div>
                <h2
                  class="font-medium text-slate-900 transition-colors"
                  :class="notification.is_read ? 'text-slate-600' : 'text-slate-900 font-semibold'"
                >
                  {{ notification.title }}
                </h2>
                <p class="text-sm text-slate-500 mt-1 leading-relaxed">
                  {{ notification.message }}
                </p>
              </div>

              <!-- Horodatage en haut à droite, discret -->
              <span class="text-xs text-slate-400 whitespace-nowrap pt-0.5">
                {{ new Date(notification.CreatedAt).toLocaleString(undefined, { hour: '2-digit', minute:'2-digit', day: 'numeric', month: 'short' }) }}
              </span>
            </div>

            <!-- Actions contextuelles en bas -->
            <div class="flex items-center gap-4 mt-4 pt-3 border-t border-slate-50">
              <button
                v-if="notification.incident_id"
                @click="goToIncident(notification)"
                class="inline-flex items-center gap-1.5 text-xs font-medium text-blue-600 hover:text-blue-700 bg-blue-50/50 hover:bg-blue-50 px-3 py-1.5 rounded-lg transition-colors"
              >
                <span>Ouvrir l'incident</span>
                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" />
                </svg>
              </button>

              <button
                v-if="!notification.is_read"
                @click="markAsRead(notification)"
                class="text-xs font-medium text-slate-400 hover:text-slate-600 px-2 py-1.5 rounded-lg hover:bg-slate-50 transition-colors"
              >
                Marquer comme lu
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- State Vide (Plus élégant) -->
      <div
        v-if="notifications.length === 0"
        class="bg-white rounded-xl border border-dashed border-slate-200 p-12 text-center max-w-4xl"
      >
        <div class="w-12 h-12 rounded-full bg-slate-50 flex items-center justify-center mx-auto mb-3">
          <svg class="w-6 h-6 text-slate-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0" />
          </svg>
        </div>
        <h3 class="text-sm font-medium text-slate-900">Aucune notification</h3>
        <p class="text-xs text-slate-400 mt-1">Vous êtes à jour ! Tout fonctionne correctement.</p>
      </div>
    </div>
  </DashboardLayout>
</template>