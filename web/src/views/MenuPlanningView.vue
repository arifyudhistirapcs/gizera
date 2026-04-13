<template>
  <div>
    <!-- Header Actions -->
    <div class="menu-planning-header">
      <a-space :size="12">
        <a-button @click="duplicatePreviousWeek">
          <template #icon><CopyOutlined /></template>
          Duplikat Minggu Lalu
        </a-button>
        <a-button type="primary" @click="showCreateModal">
          <template #icon><PlusOutlined /></template>
          Buat Menu Baru
        </a-button>
      </a-space>
    </div>

    <!-- Week Selector Card -->
    <div class="h-card week-selector">
      <a-row :gutter="16" align="middle">
        <a-col :xs="24" :md="12">
          <a-space>
            <a-button @click="previousWeek">
              <template #icon><LeftOutlined /></template>
            </a-button>
            <a-date-picker
              v-model:value="selectedWeekStart"
              picker="week"
              format="[Minggu] w, YYYY"
              @change="onWeekChange"
              style="width: 250px"
            />
            <a-button @click="nextWeek">
              <template #icon><RightOutlined /></template>
            </a-button>
            <a-button type="link" @click="goToCurrentWeek">
              Minggu Ini
            </a-button>
          </a-space>
        </a-col>
        <a-col :xs="24" :md="12" style="text-align: right">
          <a-space v-if="currentMenuPlan">
            <a-tag :color="currentMenuPlan.status === 'approved' ? 'green' : 'orange'">
              {{ currentMenuPlan.status === 'approved' ? 'Disetujui' : 'Draft' }}
            </a-tag>
            <a-button
              v-if="currentMenuPlan.status === 'draft' && canApprove"
              type="primary"
              @click="approveMenu"
              :loading="approving"
            >
              <template #icon><CheckOutlined /></template>
              Setujui Menu
            </a-button>
          </a-space>
        </a-col>
      </a-row>
    </div>

    <!-- Weekly Calendar -->
    <a-spin :spinning="loading">
      <div class="weekly-calendar">
        <a-row :gutter="[16, 16]">
          <a-col
            v-for="day in weekDays"
            :key="day.date"
            :xs="24"
            :sm="12"
            :md="12"
            :lg="6"
          >
            <div class="h-card day-card" :class="{ 'today': isToday(day.date) }">
              <!-- Day Header -->
              <div class="day-header">
                <div class="day-name">{{ day.dayName }}</div>
                <div class="day-date">{{ formatDate(day.date) }}</div>
              </div>

              <!-- Menu Items for this day -->
              <div
                class="menu-items-container"
                @drop="onDrop($event, day.date)"
                @dragover.prevent
                @dragenter.prevent
              >
                <div
                  v-for="item in getMenuItemsForDay(day.date)"
                  :key="item.id"
                  class="menu-item h-card-hover"
                  draggable="true"
                  @dragstart="onDragStart($event, item)"
                  @click="showDetailModal(item)"
                >
                  <div v-if="item.recipe?.photo_url" class="menu-item-photo">
                    <img :src="item.recipe.photo_url" :alt="item.recipe.name" />
                  </div>
                  <div class="menu-item-content">
                    <div class="menu-item-name">{{ item.recipe?.name }}</div>
                    <div class="menu-item-portions">{{ item.portions }} porsi</div>
                    <!-- Portion summary hidden -->
                    <!--
                    <div v-if="item.school_allocations && item.school_allocations.length > 0" class="menu-item-portion-summary">
                      <template v-if="getTotalSmallPortions(item.school_allocations) > 0">
                        <span class="portion-badge">Kecil: {{ getTotalSmallPortions(item.school_allocations) }}</span>
                      </template>
                      <template v-if="getTotalLargePortions(item.school_allocations) > 0">
                        <span class="portion-badge">Besar: {{ getTotalLargePortions(item.school_allocations) }}</span>
                      </template>
                    </div>
                    -->
                    <!-- School allocations hidden for cleaner UI -->
                    <!-- <div v-if="item.school_allocations && item.school_allocations.length > 0" class="menu-item-allocations">
                      <div v-for="schoolAlloc in getGroupedAllocations(item.school_allocations)" :key="schoolAlloc.school_id" class="allocation-item">
                        <span class="school-name">{{ schoolAlloc.school_name }}</span>
                        <span class="school-portions">
                          <template v-if="schoolAlloc.category === 'SD' && schoolAlloc.portions_small > 0 && schoolAlloc.portions_large > 0">
                            <span class="portion-detail">K: {{ schoolAlloc.portions_small }}</span>
                            <span class="portion-separator">|</span>
                            <span class="portion-detail">B: {{ schoolAlloc.portions_large }}</span>
                          </template>
                          <template v-else-if="schoolAlloc.category === 'SD' && schoolAlloc.portions_small > 0">
                            <span class="portion-detail">K: {{ schoolAlloc.portions_small }}</span>
                          </template>
                          <template v-else-if="schoolAlloc.category === 'SD' && schoolAlloc.portions_large > 0">
                            <span class="portion-detail">B: {{ schoolAlloc.portions_large }}</span>
                          </template>
                          <template v-else>
                            <span class="portion-detail">B: {{ schoolAlloc.portions_large || schoolAlloc.portions_small }}</span>
                          </template>
                        </span>
                      </div>
                    </div> -->
                    <!-- <div v-else class="menu-item-allocations no-allocations">
                      <span class="no-allocation-text">Belum ada alokasi</span>
                    </div> -->
                  </div>
                  <div class="menu-item-actions">
                    <a-button
                      type="text"
                      size="small"
                      @click="showEditMenuModal(item)"
                    >
                      <template #icon><EditOutlined /></template>
                    </a-button>
                    <a-button
                      type="text"
                      size="small"
                      danger
                      @click="removeMenuItem(item)"
                    >
                      <template #icon><DeleteOutlined /></template>
                    </a-button>
                  </div>
                </div>

                <!-- Add Menu Button -->
                <a-button
                  type="dashed"
                  size="small"
                  block
                  @click="showAddMenuModal(day.date)"
                  class="add-menu-button"
                >
                  <template #icon><PlusOutlined /></template>
                  Tambah Menu
                </a-button>
              </div>

              <!-- Nutrition summary hidden for cleaner UI -->
              <!--
              <div class="nutrition-summary">
                <div class="nutrition-item">
                  <span class="label">Kalori:</span>
                  <span :class="['value', getDailyNutritionStatus(day.date, 'calories')]">
                    {{ getDailyNutrition(day.date, 'calories') }} kkal
                  </span>
                </div>
                <div class="nutrition-item">
                  <span class="label">Protein:</span>
                  <span :class="['value', getDailyNutritionStatus(day.date, 'protein')]">
                    {{ getDailyNutrition(day.date, 'protein') }} g
                  </span>
                </div>
                <div class="validation-status">
                  <a-tag
                    :color="isDailyNutritionValid(day.date) ? 'success' : 'warning'"
                    style="margin: 0"
                  >
                    {{ isDailyNutritionValid(day.date) ? '✓ Memenuhi Standar' : '⚠ Belum Memenuhi' }}
                  </a-tag>
                </div>
              </div>
              -->
            </div>
          </a-col>
        </a-row>
      </div>
    </a-spin>

    <!-- Add/Edit Menu Item Modal -->
    <a-modal
      v-model:visible="addMenuModalVisible"
      :title="editingMenuItem ? 'Edit Menu' : 'Tambah Menu'"
      @ok="editingMenuItem ? updateMenuItem() : addMenuItem()"
      :ok-text="editingMenuItem ? 'Simpan' : 'Tambah'"
      cancel-text="Batal"
      :ok-button-props="{ disabled: !isAllocationValid }"
      width="700px"
      :body-style="{ paddingTop: '24px', paddingBottom: '24px', maxHeight: '70vh', overflowY: 'auto' }"
    >
      <a-form layout="vertical" style="padding: 0 4px;">
        <a-form-item label="Pilih Resep" style="margin-bottom: 24px;">
          <a-select
            v-model:value="selectedRecipeId"
            show-search
            placeholder="Cari dan pilih resep"
            :filter-option="filterRecipeOption"
            style="width: 100%"
          >
            <a-select-option
              v-for="recipe in availableRecipes"
              :key="recipe.id"
              :value="recipe.id"
            >
              {{ recipe.name }} ({{ recipe.category }})
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="Jumlah Porsi Total" style="margin-bottom: 32px;">
          <div style="padding: 8px 12px; background: #f5f5f5; border: 1px solid #d9d9d9; border-radius: 4px;">
            <span style="font-size: 24px; font-weight: 600; color: #1890ff;">{{ selectedPortions }}</span>
            <span style="margin-left: 8px; color: #8c8c8c;">porsi</span>
          </div>
          <div style="margin-top: 8px; color: #8c8c8c; font-size: 12px;">
            Total porsi dihitung otomatis dari jumlah alokasi ke semua sekolah
          </div>
        </a-form-item>
        <a-divider style="margin: 24px 0;" />
        <a-form-item label="Alokasi Sekolah" style="margin-bottom: 0;">
          <SchoolAllocationInput
            :key="allocationComponentKey"
            v-model="schoolAllocations"
            :schools="schools"
            :total-portions="selectedPortions"
            @validation-change="handleValidationChange"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Detail Menu Modal -->
    <a-modal
      v-model:visible="detailModalVisible"
      :title="detailItem?.recipe?.name || 'Detail Menu'"
      :footer="null"
      width="500px"
    >
      <div v-if="detailItem" style="display:flex;flex-direction:column;gap:16px;">
        <div v-if="detailItem.recipe?.photo_url" style="text-align:center;">
          <img :src="detailItem.recipe.photo_url" :alt="detailItem.recipe?.name" style="width:200px;height:200px;border-radius:16px;object-fit:cover;margin:0 auto;display:block;" />
        </div>
        <div style="text-align:center;">
          <div style="font-size:22px;font-weight:700;color:#303030;">{{ detailItem.portions }} porsi</div>
          <div style="font-size:13px;color:#6B6B6B;margin-top:4px;">{{ dayjs(detailItem.date).format('dddd, DD MMMM YYYY') }}</div>
        </div>
        <div v-if="detailItem.school_allocations && detailItem.school_allocations.length > 0">
          <div style="font-size:14px;font-weight:700;color:#303030;margin-bottom:10px;">📍 Alokasi per Sekolah</div>
          <div style="display:flex;flex-direction:column;gap:8px;">
            <div
              v-for="schoolAlloc in getGroupedAllocations(detailItem.school_allocations)"
              :key="schoolAlloc.school_id"
              style="display:flex;justify-content:space-between;align-items:center;padding:10px 14px;background:#F7F8FA;border-radius:10px;"
            >
              <span style="font-size:13px;font-weight:500;color:#303030;">{{ schoolAlloc.school_name }}</span>
              <span style="font-size:13px;font-weight:700;color:#303030;">
                <template v-if="schoolAlloc.portions_small > 0">
                  <span style="background:#FDEAE7;color:#F82C17;padding:2px 8px;border-radius:4px;font-size:11px;margin-right:4px;">K: {{ schoolAlloc.portions_small }}</span>
                </template>
                <template v-if="schoolAlloc.portions_large > 0">
                  <span style="background:#D1FAE5;color:#764AF1;padding:2px 8px;border-radius:4px;font-size:11px;">B: {{ schoolAlloc.portions_large }}</span>
                </template>
              </span>
            </div>
          </div>
        </div>
        <div v-else style="text-align:center;color:#6B6B6B;font-size:13px;padding:16px 0;">
          Belum ada alokasi sekolah
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, h } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  PlusOutlined,
  CopyOutlined,
  LeftOutlined,
  RightOutlined,
  CheckOutlined,
  DeleteOutlined,
  EditOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import weekOfYear from 'dayjs/plugin/weekOfYear'
import isoWeek from 'dayjs/plugin/isoWeek'
import menuPlanningService from '@/services/menuPlanningService'
import recipeService from '@/services/recipeService'
import schoolService from '@/services/schoolService'
import { useAuthStore } from '@/stores/auth'
import SchoolAllocationInput from '@/components/SchoolAllocationInput.vue'

dayjs.extend(weekOfYear)
dayjs.extend(isoWeek)

const authStore = useAuthStore()
const canApprove = computed(() => {
  return authStore.user?.role === 'ahli_gizi' || authStore.user?.role === 'kepala_sppg'
})

const loading = ref(false)
const approving = ref(false)
const selectedWeekStart = ref(dayjs().startOf('isoWeek'))
const currentMenuPlan = ref(null)
const menuItems = ref([])
const availableRecipes = ref([])
const schools = ref([])
const addMenuModalVisible = ref(false)
const selectedDate = ref(null)
const selectedRecipeId = ref(null)
const selectedPortions = ref(100)
const schoolAllocations = ref({})
const isAllocationValid = ref(false)
const draggedItem = ref(null)
const editingMenuItem = ref(null)
const allocationComponentKey = ref(0) // Key to force re-render
const detailModalVisible = ref(false)
const detailItem = ref(null)

const showDetailModal = (item) => {
  detailItem.value = item
  detailModalVisible.value = true
}

// Minimum nutrition standards per portion
const MIN_CALORIES_PER_PORTION = 600
const MIN_PROTEIN_PER_PORTION = 15

const weekDays = computed(() => {
  const days = []
  const dayNames = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu', 'Minggu']
  
  for (let i = 0; i < 7; i++) {
    const date = selectedWeekStart.value.add(i, 'day')
    days.push({
      date: date.format('YYYY-MM-DD'),
      dayName: dayNames[i],
      dayjs: date
    })
  }
  
  return days
})

const isToday = (date) => {
  return dayjs(date).isSame(dayjs(), 'day')
}

const formatDate = (date) => {
  return dayjs(date).format('DD/MM')
}

const getMenuItemsForDay = (date) => {
  return menuItems.value.filter(item => {
    // Handle both ISO format and YYYY-MM-DD format
    const itemDate = dayjs(item.date).format('YYYY-MM-DD')
    return itemDate === date
  })
}

const getSchoolName = (schoolId) => {
  const school = schools.value.find(s => s.id === schoolId)
  return school ? school.name : `School ${schoolId}`
}

const getGroupedAllocations = (allocations) => {
  // Group allocations by school_id and combine small/large portions
  const grouped = {}
  
  allocations.forEach(alloc => {
    if (!grouped[alloc.school_id]) {
      const school = schools.value.find(s => s.id === alloc.school_id)
      grouped[alloc.school_id] = {
        school_id: alloc.school_id,
        school_name: alloc.school_name || getSchoolName(alloc.school_id),
        category: school?.category || 'SMP',
        portions_small: 0,
        portions_large: 0
      }
    }
    
    // Add portions based on portion_size field
    if (alloc.portion_size === 'small') {
      grouped[alloc.school_id].portions_small = alloc.portions
    } else if (alloc.portion_size === 'large') {
      grouped[alloc.school_id].portions_large = alloc.portions
    } else {
      // Fallback: if no portion_size field, treat as large
      grouped[alloc.school_id].portions_large = alloc.portions
    }
  })
  
  return Object.values(grouped)
}

const getTotalSmallPortions = (allocations) => {
  let total = 0
  allocations.forEach(alloc => {
    if (alloc.portion_size === 'small') {
      total += alloc.portions
    }
  })
  return total
}

const getTotalLargePortions = (allocations) => {
  let total = 0
  allocations.forEach(alloc => {
    if (alloc.portion_size === 'large') {
      total += alloc.portions
    } else if (!alloc.portion_size) {
      // Fallback: if no portion_size field, treat as large
      total += alloc.portions
    }
  })
  return total
}

const getDailyNutrition = (date, type) => {
  const items = getMenuItemsForDay(date)
  let total = 0
  
  items.forEach(item => {
    const recipe = item.recipe
    if (!recipe) return
    
    // Nutrition is per menu, multiply by portions directly
    const portionFactor = item.portions
    
    switch (type) {
      case 'calories':
        total += (recipe.total_calories || 0) * portionFactor
        break
      case 'protein':
        total += (recipe.total_protein || 0) * portionFactor
        break
      case 'carbs':
        total += (recipe.total_carbs || 0) * portionFactor
        break
      case 'fat':
        total += (recipe.total_fat || 0) * portionFactor
        break
    }
  })
  
  return total.toFixed(type === 'calories' ? 0 : 1)
}

const getDailyNutritionStatus = (date, type) => {
  const value = parseFloat(getDailyNutrition(date, type))
  
  if (type === 'calories') {
    return value >= MIN_CALORIES_PER_PORTION ? 'valid' : 'invalid'
  } else if (type === 'protein') {
    return value >= MIN_PROTEIN_PER_PORTION ? 'valid' : 'invalid'
  }
  
  return ''
}

const isDailyNutritionValid = (date) => {
  const calories = parseFloat(getDailyNutrition(date, 'calories'))
  const protein = parseFloat(getDailyNutrition(date, 'protein'))
  
  return calories >= MIN_CALORIES_PER_PORTION && protein >= MIN_PROTEIN_PER_PORTION
}

const loadMenuPlan = async () => {
  loading.value = true
  try {
    const response = await menuPlanningService.getMenuPlans()
    
    console.log('Load menu plans response:', response.data)
    
    if (response.data.menu_plans && response.data.menu_plans.length > 0) {
      // Find menu plan for current week
      const weekStart = selectedWeekStart.value.format('YYYY-MM-DD')
      const plan = response.data.menu_plans.find(p => {
        const planStart = dayjs(p.week_start).format('YYYY-MM-DD')
        return planStart === weekStart
      })
      
      if (plan) {
        console.log('Found menu plan for current week:', plan)
        console.log('Menu items:', plan.menu_items)
        currentMenuPlan.value = plan
        menuItems.value = plan.menu_items || []
      } else {
        currentMenuPlan.value = null
        menuItems.value = []
      }
    } else {
      currentMenuPlan.value = null
      menuItems.value = []
    }
  } catch (error) {
    console.error('Error loading menu plan:', error)
    currentMenuPlan.value = null
    menuItems.value = []
  } finally {
    loading.value = false
  }
}

const loadRecipes = async () => {
  try {
    const response = await recipeService.getRecipes({ is_active: true })
    availableRecipes.value = response.data.recipes || []
  } catch (error) {
    message.error('Gagal memuat data resep')
    console.error('Error loading recipes:', error)
  }
}

const loadSchools = async () => {
  try {
    const response = await schoolService.getSchools({ active_only: true })
    schools.value = response.data.schools || []
  } catch (error) {
    message.error('Gagal memuat data sekolah')
    console.error('Error loading schools:', error)
  }
}

const showCreateModal = async () => {
  if (currentMenuPlan.value) {
    message.warning('Menu untuk minggu ini sudah ada')
    return
  }
  
  try {
    const weekStart = selectedWeekStart.value.format('YYYY-MM-DD')
    
    // Create empty menu plan (menu items will be added later)
    const response = await menuPlanningService.createMenuPlan({
      week_start: weekStart,
      menu_items: [] // Empty array - items will be added via "Tambah Menu" button
    })
    
    currentMenuPlan.value = response.data.menu_plan
    message.success('Menu mingguan baru berhasil dibuat')
    loadMenuPlan()
  } catch (error) {
    message.error('Gagal membuat menu mingguan')
    console.error('Error creating menu plan:', error)
  }
}

const showAddMenuModal = (date) => {
  if (!currentMenuPlan.value) {
    message.warning('Buat menu mingguan terlebih dahulu')
    return
  }
  
  if (currentMenuPlan.value.status === 'approved') {
    message.warning('Menu yang sudah disetujui tidak dapat diubah')
    return
  }
  
  // Reset all form state for new menu
  selectedDate.value = date
  selectedRecipeId.value = null
  selectedPortions.value = 0  // Set to 0 instead of 100
  editingMenuItem.value = null
  
  // Reset school allocations to empty object
  schoolAllocations.value = {}
  
  // Reset validation
  isAllocationValid.value = false
  
  // Force re-render of SchoolAllocationInput component
  allocationComponentKey.value++
  
  // Open modal
  addMenuModalVisible.value = true
}

const showEditMenuModal = (item) => {
  if (!currentMenuPlan.value) {
    message.warning('Buat menu mingguan terlebih dahulu')
    return
  }
  
  if (currentMenuPlan.value.status === 'approved') {
    message.warning('Menu yang sudah disetujui tidak dapat diubah')
    return
  }
  
  editingMenuItem.value = item
  selectedDate.value = item.date
  selectedRecipeId.value = item.recipe_id
  selectedPortions.value = item.portions
  
  // Load existing allocations with portion sizes
  let allocations = {}
  if (item.school_allocations && item.school_allocations.length > 0) {
    // Group allocations by school_id and combine small/large portions
    const groupedAllocations = {}
    item.school_allocations.forEach(alloc => {
      if (!groupedAllocations[alloc.school_id]) {
        groupedAllocations[alloc.school_id] = {
          portions_small: 0,
          portions_large: 0
        }
      }
      
      // Add portions based on portion_size field
      if (alloc.portion_size === 'small') {
        groupedAllocations[alloc.school_id].portions_small = alloc.portions
      } else if (alloc.portion_size === 'large') {
        groupedAllocations[alloc.school_id].portions_large = alloc.portions
      } else {
        // Fallback: if no portion_size field, treat as large
        groupedAllocations[alloc.school_id].portions_large = alloc.portions
      }
    })
    
    allocations = groupedAllocations
  }
  schoolAllocations.value = allocations
  
  // Trigger validation
  isAllocationValid.value = validateAllocations()
  
  // Force re-render of SchoolAllocationInput component
  allocationComponentKey.value++
  
  addMenuModalVisible.value = true
}

const filterRecipeOption = (input, option) => {
  return option.children[0].children.toLowerCase().includes(input.toLowerCase())
}

const handleValidationChange = (validation) => {
  console.log('handleValidationChange called:', validation)
  isAllocationValid.value = validation.isValid
  // Auto-calculate total portions from allocations
  selectedPortions.value = validation.totalAllocated || 0
  console.log('selectedPortions updated to:', selectedPortions.value)
}

const validateAllocations = () => {
  let totalAllocated = 0
  Object.values(schoolAllocations.value).forEach(alloc => {
    if (alloc && typeof alloc === 'object') {
      totalAllocated += (alloc.portions_small || 0) + (alloc.portions_large || 0)
    }
  })
  
  // Since selectedPortions is auto-calculated from allocations,
  // we only need to check if there are any allocations
  if (totalAllocated === 0) return false
  return true
}

const handlePortionsChange = () => {
  // Reset allocations when portions change
  schoolAllocations.value = {}
  isAllocationValid.value = false
}

const addMenuItem = async () => {
  console.log('addMenuItem called')
  console.log('selectedRecipeId:', selectedRecipeId.value)
  console.log('selectedPortions:', selectedPortions.value)
  console.log('isAllocationValid:', isAllocationValid.value)
  console.log('schoolAllocations:', schoolAllocations.value)
  
  if (!selectedRecipeId.value) {
    message.warning('Pilih resep terlebih dahulu')
    return
  }
  
  if (!isAllocationValid.value) {
    message.warning('Alokasi sekolah belum valid')
    return
  }
  
  try {
    // Transform allocations to API format with portion sizes
    const school_allocations = Object.entries(schoolAllocations.value)
      .filter(([_, alloc]) => {
        // Include if either portion type > 0
        return (alloc.portions_small > 0 || alloc.portions_large > 0)
      })
      .map(([school_id, alloc]) => ({
        school_id: parseInt(school_id),
        portions_small: alloc.portions_small || 0,
        portions_large: alloc.portions_large || 0
      }))
    
    const payload = {
      date: selectedDate.value,
      recipe_id: selectedRecipeId.value,
      portions: selectedPortions.value,
      school_allocations
    }
    
    console.log('Creating menu item with payload:', payload)
    
    // Call the new createMenuItem endpoint
    const response = await menuPlanningService.createMenuItem(currentMenuPlan.value.id, payload)
    
    console.log('Create menu item response:', response.data)
    
    addMenuModalVisible.value = false
    message.success('Menu berhasil ditambahkan')
    
    // Reload menu plan to get updated data
    await loadMenuPlan()
  } catch (error) {
    message.error('Gagal menambahkan menu')
    console.error('Error adding menu item:', error)
    console.error('Error response:', error.response?.data)
  }
}

const updateMenuItem = async () => {
  if (!selectedRecipeId.value) {
    message.warning('Pilih resep terlebih dahulu')
    return
  }
  
  if (!isAllocationValid.value) {
    message.warning('Alokasi sekolah belum valid')
    return
  }
  
  try {
    // Transform allocations to API format with portion sizes
    const school_allocations = Object.entries(schoolAllocations.value)
      .filter(([_, alloc]) => {
        // Include if either portion type > 0
        return (alloc.portions_small > 0 || alloc.portions_large > 0)
      })
      .map(([school_id, alloc]) => ({
        school_id: parseInt(school_id),
        portions_small: alloc.portions_small || 0,
        portions_large: alloc.portions_large || 0
      }))
    
    const payload = {
      date: selectedDate.value,
      recipe_id: selectedRecipeId.value,
      portions: selectedPortions.value,
      school_allocations
    }
    
    // Call the new updateMenuItem endpoint
    await menuPlanningService.updateMenuItem(
      currentMenuPlan.value.id,
      editingMenuItem.value.id,
      payload
    )
    
    addMenuModalVisible.value = false
    editingMenuItem.value = null
    message.success('Menu berhasil diperbarui')
    
    // Reload menu plan to get updated data
    await loadMenuPlan()
  } catch (error) {
    message.error('Gagal memperbarui menu')
    console.error('Error updating menu item:', error)
  }
}

const removeMenuItem = async (item) => {
  if (currentMenuPlan.value.status === 'approved') {
    message.warning('Menu yang sudah disetujui tidak dapat diubah')
    return
  }
  
  try {
    // Call the delete endpoint directly
    await menuPlanningService.deleteMenuItem(currentMenuPlan.value.id, item.id)
    message.success('Menu berhasil dihapus')
    
    // Close detail modal
    detailModalVisible.value = false
    detailItem.value = null
    
    // Reload menu plan to get updated data
    await loadMenuPlan()
  } catch (error) {
    message.error('Gagal menghapus menu')
    console.error('Error removing menu item:', error)
  }
}

const saveMenuPlan = async (items) => {
  try {
    const payload = {
      week_start: currentMenuPlan.value.week_start,
      week_end: currentMenuPlan.value.week_end,
      status: currentMenuPlan.value.status,
      menu_items: items.map(item => ({
        date: item.date,
        recipe_id: item.recipe_id,
        portions: item.portions,
        school_allocations: item.school_allocations || []
      }))
    }
    
    await menuPlanningService.updateMenuPlan(currentMenuPlan.value.id, payload)
    // Reload menu plan to get updated data
    await loadMenuPlan()
  } catch (error) {
    throw error
  }
}

const approveMenu = async () => {
  // Check for empty days and insufficient nutrition
  const emptyDays = []
  const insufficientDays = []
  
  weekDays.value.forEach(day => {
    const items = getMenuItemsForDay(day.date)
    if (items.length === 0) {
      emptyDays.push(day.dayName)
    } else if (!isDailyNutritionValid(day.date)) {
      const calories = parseFloat(getDailyNutrition(day.date, 'calories'))
      const protein = parseFloat(getDailyNutrition(day.date, 'protein'))
      insufficientDays.push({
        day: day.dayName,
        calories: calories.toFixed(0),
        protein: protein.toFixed(1)
      })
    }
  })
  
  // Show confirmation modal if there are issues
  if (emptyDays.length > 0 || insufficientDays.length > 0) {
    let content = ''
    
    if (emptyDays.length > 0) {
      content += `<p><strong>Hari yang belum diisi menu:</strong></p>`
      content += `<ul>${emptyDays.map(d => `<li>${d}</li>`).join('')}</ul>`
    }
    
    if (insufficientDays.length > 0) {
      content += `<p><strong>Hari dengan nutrisi tidak memenuhi standar (min 600 kcal, 15g protein):</strong></p>`
      content += `<ul>${insufficientDays.map(d => `<li>${d.day}: ${d.calories} kcal, ${d.protein}g protein</li>`).join('')}</ul>`
    }
    
    content += `<p style="margin-top: 16px;">Apakah Anda yakin ingin menyetujui menu ini?</p>`
    
    Modal.confirm({
      title: 'Konfirmasi Persetujuan Menu',
      content: h('div', { innerHTML: content }),
      okText: 'Ya, Setujui',
      cancelText: 'Batal',
      onOk: async () => {
        await performApprove()
      }
    })
  } else {
    // No issues, approve directly
    await performApprove()
  }
}

const performApprove = async () => {
  approving.value = true
  try {
    await menuPlanningService.approveMenuPlan(currentMenuPlan.value.id)
    message.success('Menu berhasil disetujui')
    await loadMenuPlan()
  } catch (error) {
    message.error('Gagal menyetujui menu')
    console.error('Error approving menu:', error)
  } finally {
    approving.value = false
  }
}

const duplicatePreviousWeek = async () => {
  if (currentMenuPlan.value) {
    message.warning('Menu untuk minggu ini sudah ada')
    return
  }
  
  try {
    const response = await menuPlanningService.getMenuPlans()
    
    if (!response.data.menu_plans || response.data.menu_plans.length === 0) {
      message.warning('Tidak ada menu minggu lalu untuk diduplikat')
      return
    }
    
    // Find menu from previous week
    const previousWeekStart = selectedWeekStart.value.subtract(7, 'day').format('YYYY-MM-DD')
    const previousMenu = response.data.menu_plans.find(p => {
      const planStart = dayjs(p.week_start).format('YYYY-MM-DD')
      return planStart === previousWeekStart
    })
    
    if (!previousMenu) {
      message.warning('Tidak ada menu minggu lalu untuk diduplikat')
      return
    }
    
    const weekStart = selectedWeekStart.value.format('YYYY-MM-DD')
    
    // Create new menu with items from previous week
    const newMenuResponse = await menuPlanningService.createMenuPlan({
      week_start: weekStart,
      menu_items: previousMenu.menu_items?.map(item => {
        // Calculate date offset
        const oldDate = dayjs(item.date)
        const newDate = oldDate.add(7, 'day').format('YYYY-MM-DD')
        
        return {
          date: newDate,
          recipe_id: item.recipe_id,
          portions: item.portions
        }
      }) || []
    })
    
    currentMenuPlan.value = newMenuResponse.data.menu_plan
    menuItems.value = newMenuResponse.data.menu_plan.menu_items || []
    message.success('Menu minggu lalu berhasil diduplikat')
  } catch (error) {
    message.error('Gagal menduplikat menu minggu lalu')
    console.error('Error duplicating previous week:', error)
  }
}

const onDragStart = (event, item) => {
  draggedItem.value = item
  event.dataTransfer.effectAllowed = 'move'
}

const onDrop = async (event, targetDate) => {
  if (!draggedItem.value) return
  
  if (currentMenuPlan.value.status === 'approved') {
    message.warning('Menu yang sudah disetujui tidak dapat diubah')
    return
  }
  
  try {
    // Update the date of the dragged item
    const updatedItems = menuItems.value.map(item => {
      if (item === draggedItem.value) {
        return { ...item, date: targetDate }
      }
      return item
    })
    
    await saveMenuPlan(updatedItems)
    message.success('Menu berhasil dipindahkan')
  } catch (error) {
    message.error('Gagal memindahkan menu')
    console.error('Error moving menu item:', error)
  } finally {
    draggedItem.value = null
  }
}

const onWeekChange = () => {
  loadMenuPlan()
}

const previousWeek = () => {
  selectedWeekStart.value = selectedWeekStart.value.subtract(7, 'day')
  loadMenuPlan()
}

const nextWeek = () => {
  selectedWeekStart.value = selectedWeekStart.value.add(7, 'day')
  loadMenuPlan()
}

const goToCurrentWeek = () => {
  selectedWeekStart.value = dayjs().startOf('isoWeek')
  loadMenuPlan()
}

onMounted(() => {
  loadMenuPlan()
  loadRecipes()
  loadSchools()
})
</script>

<style scoped>
/* Header Actions */
.menu-planning-header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: 16px;
}

/* Week Selector */
.week-selector {
  padding: 16px 20px;
  margin-bottom: 16px;
  border-radius: 14px;
}

/* Day Card — Modern, proportional */
.day-card {
  min-height: 220px;
  display: flex;
  flex-direction: column;
  padding: 16px;
  border-radius: 14px !important;
  border: 1px solid #F0F0F0 !important;
  transition: all 0.2s ease;
}

.day-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
}

.day-card.today {
  border: 2px solid #F82C17 !important;
  box-shadow: 0 0 0 3px rgba(248, 44, 23, 0.1);
}

/* Day Header — Compact */
.day-header {
  text-align: center;
  padding-bottom: 10px;
  border-bottom: 1px solid #F0F0F0;
  margin-bottom: 12px;
}

.day-name {
  font-size: 16px;
  font-weight: 700;
  color: #303030;
  margin-bottom: 2px;
}

.day-date {
  font-size: 13px;
  font-weight: 500;
  color: #6B6B6B;
}

/* Menu Items Container */
.menu-items-container {
  flex: 1;
  min-height: 100px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 12px;
}

/* Menu Item — Compact, clean */
.menu-item {
  background: #FAFAFA;
  padding: 10px;
  border-radius: 10px;
  cursor: move;
  display: flex;
  gap: 10px;
  align-items: flex-start;
  transition: all 0.2s ease;
  border: 1px solid #F0F0F0;
}

.menu-item:hover {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transform: translateY(-1px);
}

.menu-item-photo {
  flex-shrink: 0;
  width: 52px;
  height: 52px;
  border-radius: 8px;
  overflow: hidden;
  background: #F0F0F0;
}

.menu-item-photo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.menu-item-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
}

.menu-item-actions {
  display: flex;
  gap: 2px;
  flex-shrink: 0;
  opacity: 0.5;
  transition: opacity 0.2s;
}

.menu-item:hover .menu-item-actions {
  opacity: 1;
}

.menu-item-name {
  font-weight: 600;
  font-size: 13px;
  color: #303030;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.menu-item-portions {
  font-size: 11px;
  color: #6B6B6B;
  font-weight: 500;
}

.menu-item-portion-summary {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.portion-badge {
  padding: 1px 6px;
  background: #FDEAE7;
  border: none;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 600;
  color: #F82C17;
}

/* School Allocations — Compact */
.menu-item-allocations {
  display: flex;
  flex-direction: column;
  gap: 2px;
  max-height: 80px;
  overflow-y: auto;
  scrollbar-width: thin;
}

.allocation-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
  padding: 2px 6px;
  background: #F7F8FA;
  border-radius: 4px;
}

.allocation-item .school-name {
  color: #6B6B6B;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 80px;
}

.allocation-item .school-portions {
  color: #303030;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 2px;
  font-size: 10px;
}

.portion-detail { white-space: nowrap; }
.portion-separator { color: #D8D8DB; margin: 0 1px; }

.no-allocations {
  font-size: 11px;
  color: #F82C17;
  font-style: italic;
}

.no-allocation-text {
  padding: 2px 6px;
}

/* Add Menu Button */
.add-menu-button {
  margin-top: auto;
  height: 36px;
  border-radius: 8px;
  border-color: #D8D8DB;
  color: #6B6B6B;
  font-weight: 500;
  font-size: 13px;
}

.add-menu-button:hover {
  border-color: #F82C17;
  color: #F82C17;
  background: rgba(248, 44, 23, 0.04);
}

/* Nutrition Summary — Compact */
.nutrition-summary {
  padding-top: 10px;
  border-top: 1px solid #F0F0F0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nutrition-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.nutrition-item .label {
  color: #6B6B6B;
  font-weight: 500;
}

.nutrition-item .value {
  font-weight: 600;
  color: #303030;
}

.nutrition-item .value.valid {
  color: #764AF1;
}

.nutrition-item .value.invalid {
  color: #D97706;
}

.validation-status {
  margin-top: 4px;
  text-align: center;
}

/* Responsive */
@media (max-width: 1024px) {
  .day-card { min-height: 300px; }
}

@media (max-width: 767px) {
  .menu-planning-header {
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
  }
  .week-selector :deep(.ant-row) { flex-direction: column; gap: 10px; }
  .week-selector :deep(.ant-col) { text-align: left !important; }
  .day-card { min-height: 280px; }
  .menu-item-photo { width: 44px; height: 44px; }
}

/* Dark Mode */
.dark .day-card { background: #252525; border-color: #404040 !important; }
.dark .day-card.today { border-color: #F82C17 !important; }
.dark .day-header { border-bottom-color: #404040; }
.dark .day-name { color: #F7F8FA; }
.dark .menu-item { background: #303030; border-color: #404040; }
.dark .menu-item:hover { background: #353535; }
.dark .menu-item-name { color: #F7F8FA; }
.dark .allocation-item { background: #353535; }
.dark .nutrition-summary { border-top-color: #404040; }
</style>
