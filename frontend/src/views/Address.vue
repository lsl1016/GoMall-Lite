<script setup>
import { onMounted, reactive, ref } from 'vue'
import NavBar from '@/components/NavBar.vue'
import EmptyState from '@/components/EmptyState.vue'
import { addAddressAPI, deleteAddressAPI, getAddressesAPI, setDefaultAddressAPI, updateAddressAPI } from '@/api/address'

const addresses = ref([])
const loading = ref(false)
const errorMsg = ref('')
const showModal = ref(false)
const editingId = ref(null)
const form = reactive({ receiver: '', phone: '', province: '', city: '', district: '', detail: '', isDefault: false })

function resetForm() {
  editingId.value = null
  Object.assign(form, { receiver: '', phone: '', province: '', city: '', district: '', detail: '', isDefault: false })
}

async function fetchAddresses() {
  try {
    loading.value = true
    errorMsg.value = ''
    addresses.value = await getAddressesAPI()
  } catch (error) {
    errorMsg.value = error.message
  } finally {
    loading.value = false
  }
}

function openAdd() {
  resetForm()
  showModal.value = true
}

function openEdit(address) {
  editingId.value = address.id
  Object.assign(form, address)
  showModal.value = true
}

async function save() {
  try {
    if (editingId.value) {
      addresses.value = await updateAddressAPI(editingId.value, form)
    } else {
      addresses.value = await addAddressAPI(form)
    }
    showModal.value = false
  } catch (error) {
    errorMsg.value = error.message
  }
}

async function remove(id) {
  addresses.value = await deleteAddressAPI(id)
}

async function setDefault(id) {
  addresses.value = await setDefaultAddressAPI(id)
}

onMounted(fetchAddresses)
</script>

<template>
  <div class="page">
    <NavBar />
    <main class="container main">
      <section class="card">
        <div class="address-page-head">
          <h1 style="font-size: 34px; margin: 0">地址管理</h1>
          <button class="primary-btn small" @click="openAdd">+ 新增地址</button>
        </div>
        <div v-if="loading" class="loading">加载中...</div>
        <p v-else-if="errorMsg" class="error-text" style="padding: 0 30px">{{ errorMsg }}</p>
        <EmptyState v-else-if="!addresses.length" title="暂无收货地址" button-text="新增地址" @action="openAdd" />
        <div v-else class="address-list">
          <div v-for="addr in addresses" :key="addr.id" class="address-card" :class="{ default: addr.isDefault }">
            <div>
              <strong style="font-size: 18px">{{ addr.receiver }}</strong>
              <span style="margin-left: 20px">{{ addr.phone }}</span>
              <span v-if="addr.isDefault" class="tag" style="margin-left: 12px">默认地址</span>
              <p class="muted">{{ addr.province }} {{ addr.city }} {{ addr.district }} {{ addr.detail }}</p>
            </div>
            <div class="address-actions">
              <button v-if="!addr.isDefault" class="ghost-btn small" @click="setDefault(addr.id)">设为默认</button>
              <button class="text-btn" @click="openEdit(addr)">编辑</button>
              <button class="text-danger" @click="remove(addr.id)">删除</button>
            </div>
          </div>
        </div>
      </section>
    </main>

    <div v-if="showModal" class="modal-mask">
      <div class="modal">
        <h2>{{ editingId ? '编辑地址' : '新增地址' }}</h2>
        <div class="form-grid-2">
          <div class="form-row"><label>收货人</label><input v-model="form.receiver" class="input" /></div>
          <div class="form-row"><label>手机号</label><input v-model="form.phone" class="input" /></div>
          <div class="form-row"><label>省</label><input v-model="form.province" class="input" placeholder="北京市" /></div>
          <div class="form-row"><label>市</label><input v-model="form.city" class="input" placeholder="北京市" /></div>
          <div class="form-row"><label>区</label><input v-model="form.district" class="input" placeholder="朝阳区" /></div>
          <div class="form-row"><label>默认地址</label><label style="display:flex;align-items:center;height:44px"><input v-model="form.isDefault" type="checkbox" /> 设为默认</label></div>
        </div>
        <div class="form-row"><label>详细地址</label><input v-model="form.detail" class="input" placeholder="街道门牌号" /></div>
        <div class="modal-actions">
          <button class="ghost-btn" @click="showModal = false">取消</button>
          <button class="primary-btn" @click="save">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>
