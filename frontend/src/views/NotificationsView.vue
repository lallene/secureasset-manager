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
  if (!notification.incident_id) {
    return;
  }

  if (!notification.is_read) {
    await markAsRead(notification);
  }

  router.push(`/incidents?incident=${notification.incident_id}`);
};

const fetchNotifications = async () => {
  try {
    const response = await api.get("/notifications", {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
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
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    notification.is_read = true;
  } catch (err) {
    console.error(err);
  }
};

onMounted(fetchNotifications);
</script>

<template>
  <DashboardLayout>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">Notifications</h1>

      <span class="bg-red-600 text-white px-3 py-1 rounded-full text-sm">
        {{ unreadCount }} non lues
      </span>
    </div>

    <p v-if="error" class="mb-4 text-red-600">
      {{ error }}
    </p>

    <div class="space-y-4">
      <div
        v-for="notification in notifications"
        :key="notification.ID"
        class="bg-white rounded-xl shadow p-5 border-l-4"
        :class="notification.is_read ? 'border-slate-300 opacity-70' : 'border-blue-600'"
      >
        <div class="flex items-start justify-between">
          <div>
            <h2 class="font-bold text-lg">
              {{ notification.title }}
            </h2>

            <p class="text-gray-600 mt-1">
              {{ notification.message }}
            </p>
            <button
              v-if="notification.incident_id"
              @click="goToIncident(notification)"
              class="mt-3 text-blue-600 hover:underline text-sm"
            >
              Voir le ticket
            </button>

            <p class="text-xs text-gray-400 mt-2">
              {{ new Date(notification.CreatedAt).toLocaleString() }}
            </p>
          </div>

          <button
            v-if="!notification.is_read"
            @click="markAsRead(notification)"
            class="text-blue-600 hover:underline"
          >
            Marquer comme lu
          </button>
        </div>
      </div>

      <div
        v-if="notifications.length === 0"
        class="bg-white rounded-xl shadow p-6 text-center text-gray-500"
      >
        Aucune notification.
      </div>
    </div>
  </DashboardLayout>
</template>