<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";

import api from "../api/axios";

const router = useRouter();

const email = ref("");
const password = ref("");
const error = ref("");

const login = async () => {
  try {
    const response = await api.post("/login", {
      email: email.value,
      password: password.value,
    });

    localStorage.setItem("token", response.data.token);

    router.push("/dashboard");
  } catch (err) {
    error.value = "Email ou mot de passe incorrect";
  }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-950">
    <div class="bg-white p-8 rounded-xl shadow w-full max-w-md">
      <h1 class="text-2xl font-bold mb-6 text-center">
        SecureAsset Manager
      </h1>

      <div class="flex flex-col gap-4">
        <input
          v-model="email"
          type="email"
          placeholder="Email"
          class="border rounded-lg p-3"
        />

        <input
          v-model="password"
          type="password"
          placeholder="Mot de passe"
          class="border rounded-lg p-3"
        />

        <button
          @click="login"
          class="bg-blue-600 hover:bg-blue-700 text-white rounded-lg p-3"
        >
          Connexion
        </button>

        <p v-if="error" class="text-red-500 text-sm">
          {{ error }}
        </p>
      </div>
    </div>
  </div>
</template>