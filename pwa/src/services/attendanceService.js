/**
 * Attendance Service for PWA
 * Handles check-in/out operations with Wi-Fi validation
 */

import { attendanceAPI, wifiAPI } from './api.js'
import wifiService from './wifiService.js'
import { useAuthStore } from '@/stores/auth.js'

class AttendanceService {
  constructor() {
    this.currentAttendance = null
    this.authorizedNetworks = []
  }

  /**
   * Initialize attendance service
   */
  async initialize() {
    try {
      await this.loadAuthorizedNetworks()
    } catch (error) {
      console.error('Failed to initialize attendance service:', error)
    }
  }

  /**
   * Load authorized Wi-Fi networks from backend
   */
  async loadAuthorizedNetworks() {
    try {
      const response = await wifiAPI.getAuthorizedNetworks()
      
      // Backend returns { success: true, data: [...] }
      if (response.data.success && response.data.data) {
        this.authorizedNetworks = response.data.data.map(config => ({
          ssid: config.ssid,
          bssid: config.bssid,
          location: config.location || 'Kantor',
          is_active: config.is_active
        }))
      } else {
        // Fallback to empty array if no data
        this.authorizedNetworks = []
      }
      
      console.log('Loaded authorized networks:', this.authorizedNetworks)
      return this.authorizedNetworks
    } catch (error) {
      console.error('Failed to load authorized networks:', error)
      // Fallback to empty array
      this.authorizedNetworks = []
      return this.authorizedNetworks
    }
  }

  /**
   * Get current attendance status
   */
  async getCurrentAttendance() {
    try {
      const response = await attendanceAPI.getCurrentAttendance()
      if (response.data.success && response.data.data) {
        this.currentAttendance = response.data.data
      } else {
        this.currentAttendance = null
      }
      return this.currentAttendance
    } catch (error) {
      console.error('Failed to get current attendance:', error)
      return null
    }
  }

  /**
   * Perform check-in with Wi-Fi validation
   * @param {string} manualSSID - Optional manual SSID input
   * @param {boolean} useGPS - Whether to use GPS validation
   * @returns {Promise<Object>}
   */
  async checkIn(manualSSID = null, useGPS = false) {
    try {
      const authStore = useAuthStore()
      
      // Step 1: Get GPS location if requested
      let gpsLocation = null
      if (useGPS) {
        try {
          gpsLocation = await wifiService.getCurrentLocation()
        } catch (error) {
          return {
            success: false,
            error: 'GPS Error',
            message: error.message
          }
        }
      }

      // Step 2: Validate Wi-Fi connection
      const wifiValidation = await wifiService.validateWiFiConnection(
        manualSSID,
        gpsLocation,
        this.authorizedNetworks
      )

      if (!wifiValidation.isValid) {
        return {
          success: false,
          error: 'Wi-Fi Validation Failed',
          message: wifiValidation.error,
          details: wifiValidation.details,
          method: wifiValidation.method
        }
      }

      // Step 3: Submit check-in to backend
      const checkInData = {
        employee_id: authStore.user.id,
        check_in_time: new Date().toISOString(),
        wifi_validation: wifiValidation,
        location: gpsLocation
      }

      const response = await attendanceAPI.checkIn(checkInData)
      
      this.currentAttendance = response.data.attendance
      
      return {
        success: true,
        message: 'Check-in berhasil!',
        attendance: this.currentAttendance,
        validation: wifiValidation
      }

    } catch (error) {
      console.error('Check-in failed:', error)
      return {
        success: false,
        error: 'Check-in Failed',
        message: error.response?.data?.message || 'Terjadi kesalahan saat check-in',
        details: error.message
      }
    }
  }

  /**
   * Perform check-out
   * @returns {Promise<Object>}
   */
  async checkOut() {
    try {
      if (!this.currentAttendance || this.currentAttendance.check_out || this.currentAttendance.check_out_time) {
        return {
          success: false,
          error: 'Invalid State',
          message: 'Tidak ada check-in aktif atau sudah check-out'
        }
      }

      // Check-out doesn't require WiFi/IP validation
      // Employee can check-out from anywhere
      // Send empty request body
      const response = await attendanceAPI.checkOut({})
      
      if (response.data.success) {
        this.currentAttendance = response.data.data
        
        return {
          success: true,
          message: 'Check-out berhasil!',
          attendance: this.currentAttendance,
          workHours: this.currentAttendance.work_hours
        }
      }
      
      return {
        success: false,
        error: 'Check-out Failed',
        message: response.data.message || 'Gagal melakukan check-out'
      }

    } catch (error) {
      console.error('Check-out failed:', error)
      return {
        success: false,
        error: 'Check-out Failed',
        message: error.response?.data?.message || 'Terjadi kesalahan saat check-out',
        details: error.message
      }
    }
  }

  /**
   * Get attendance history
   * @param {number} days - Number of days to fetch
   * @returns {Promise<Array>}
   */
  async getAttendanceHistory(days = 7) {
    try {
      const response = await attendanceAPI.getHistory(days)
      if (response.data.success && response.data.data) {
        return response.data.data
      }
      return []
    } catch (error) {
      console.error('Failed to get attendance history:', error)
      return []
    }
  }

  /**
   * Validate Wi-Fi without checking in (for testing)
   * @param {string} manualSSID 
   * @param {boolean} useGPS 
   * @returns {Promise<Object>}
   */
  async validateWiFiOnly(manualSSID = null, useGPS = false) {
    try {
      let gpsLocation = null
      if (useGPS) {
        gpsLocation = await wifiService.getCurrentLocation()
      }

      return await wifiService.validateWiFiConnection(
        manualSSID,
        gpsLocation,
        this.authorizedNetworks
      )
    } catch (error) {
      return {
        isValid: false,
        error: 'Validation Error',
        details: error.message
      }
    }
  }

  /**
   * Get authorized networks (for UI display)
   * @returns {Array}
   */
  getAuthorizedNetworks() {
    return this.authorizedNetworks.map(network => ({
      ssid: network.ssid,
      location: network.location || 'Kantor'
    }))
  }

  /**
   * Format work hours for display
   * @param {number} hours 
   * @returns {string}
   */
  formatWorkHours(hours) {
    if (!hours) return '0 jam 0 menit'
    
    const wholeHours = Math.floor(hours)
    const minutes = Math.round((hours - wholeHours) * 60)
    
    return `${wholeHours} jam ${minutes} menit`
  }

  /**
   * Check if user can check in
   * @returns {boolean}
   */
  canCheckIn() {
    // Can check in if:
    // 1. No current attendance, OR
    // 2. Already checked out (check_out exists)
    return !this.currentAttendance || 
           (this.currentAttendance && (this.currentAttendance.check_out || this.currentAttendance.check_out_time))
  }

  /**
   * Check if user can check out
   * @returns {boolean}
   */
  canCheckOut() {
    // Can check out if:
    // 1. Has current attendance, AND
    // 2. Has checked in (check_in exists), AND
    // 3. Has NOT checked out yet (check_out is null)
    return this.currentAttendance && 
           (this.currentAttendance.check_in || this.currentAttendance.check_in_time) && 
           !(this.currentAttendance.check_out || this.currentAttendance.check_out_time)
  }
}

export default new AttendanceService()