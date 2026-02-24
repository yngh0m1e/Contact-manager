<template>
  <div class="min-h-screen bg-gradient-to-br from-[#0f172a] via-[#1e2937] to-[#4c1d95] flex items-center justify-center">
    <div class="w-full max-w-md px-6">
      <div class="text-center mb-12">
        <div class="mx-auto w-20 h-20 bg-white/10 backdrop-blur-3xl rounded-3xl flex items-center justify-center text-6xl mb-6 border border-white/30">
          📇
        </div>
        <h1 class="text-5xl font-light text-white tracking-tighter">Contact Manager</h1>
        <p class="text-zinc-400 mt-3">Управляй контактами красиво</p>
      </div>

      <div class="bg-white/10 backdrop-blur-3xl border border-white/20 rounded-3xl p-10">
        <h2 class="text-3xl font-light text-white text-center mb-8">Добро пожаловать!</h2>

        <div class="space-y-6">
          <input v-model="username" placeholder="Логин" 
            class="w-full bg-white/10 border border-white/30 rounded-2xl px-6 py-4 text-white placeholder-zinc-400 focus:border-[#9f1239] outline-none transition" />
          <input v-model="password" type="password" placeholder="Пароль" 
            class="w-full bg-white/10 border border-white/30 rounded-2xl px-6 py-4 text-white placeholder-zinc-400 focus:border-[#9f1239] outline-none transition" />

          <button @click="login" 
            class="w-full bg-[#9f1239] hover:bg-[#b91c1c] text-white font-medium py-4 rounded-2xl transition active:scale-95">
            Войти
          </button>

          <button @click="showRegister = !showRegister" class="w-full text-zinc-400 hover:text-white transition">
            Ещё нет аккаунта? Зарегистрироваться.
          </button>

          <div v-if="showRegister" class="pt-6 border-t border-white/20 space-y-5">
            <input v-model="regData.username" placeholder="Логин" class="w-full bg-white/10 border border-white/30 rounded-2xl px-6 py-4 text-white" />
            <input v-model="regData.password" type="password" placeholder="Пароль" class="w-full bg-white/10 border border-white/30 rounded-2xl px-6 py-4 text-white" />
            <input v-model="regData.email" type="email" placeholder="Email" class="w-full bg-white/10 border border-white/30 rounded-2xl px-6 py-4 text-white" />
            <input v-model="regData.name" placeholder="Имя" class="w-full bg-white/10 border border-white/30 rounded-2xl px-6 py-4 text-white" />
            
            <button @click="register" class="w-full bg-emerald-600 hover:bg-emerald-700 py-4 rounded-2xl text-white">
              Создать аккаунт
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

const username = ref('')
const password = ref('')
const showRegister = ref(false)
const regData = ref({ username: '', password: '', email: '', name: '' })

const login = async () => {
  await auth.login(username.value, password.value)
  router.push('/contacts')
}

const register = async () => {
  await auth.register(regData.value)
  alert('Аккаунт создан!')
  showRegister.value = false
}
</script>