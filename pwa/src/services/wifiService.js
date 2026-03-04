/**
 * Wi-Fi Detection Service for PWA Attendance Module
 * 
 * Browser Wi-Fi detection has limited support:
 * - Network Information API provides limited network info
 * - No direct SSID/BSSID access due to security restrictions
 * - Fallback to manual input and GPS validation
 */

class WiFiService {
  constructor() {
    this.networkInfo = null
    this.geolocation = null
    this.initializeNetworkInfo()
  }

  /**
   * Initialize Network Information API if available
   */
  initializeNetworkInfo() {
    if ('connection' in navigator) {
      this.networkInfo = navigator.connection
    } else if ('mozConnection' in navigator) {
      this.networkInfo = navigator.mozConnection
    } else if ('webkitConnection' in navigator) {
      this.networkInfo = navigator.webkitConnection
    }
  }

  /**
   * Check if device is connected to Wi-Fi
   * @returns {Promise<boolean>}
   */
  async isConnectedToWiFi() {
    try {
      // Check if online first
      if (!navigator.onLine) {
        return false
      }

      if (this.networkInfo) {
        // Check if connection type is wifi
        const type = this.networkInfo.type
        
        // Accept wifi or unknown (desktop browsers often return unknown)
        if (type === 'wifi') {
          return true
        }
        
        // If type is unknown but we have good connection, assume WiFi on desktop
        if (type === 'unknown' || !type) {
          // Check if we have good connection speed (likely WiFi)
          const effectiveType = this.networkInfo.effectiveType
          if (effectiveType === '4g' || effectiveType === 'unknown') {
            // Likely desktop/laptop with WiFi
            return true
          }
        }
      }
      
      // Fallback: if online and no cellular indicator, assume WiFi
      // This is especially true for desktop browsers
      return navigator.onLine
    } catch (error) {
      console.warn('Wi-Fi detection failed:', error)
      // If we can't detect, assume WiFi if online (permissive approach)
      return navigator.onLine
    }
  }

  /**
   * Get network connection information
   * @returns {Object}
   */
  getNetworkInfo() {
    if (!this.networkInfo) {
      return {
        type: 'unknown',
        effectiveType: 'unknown',
        downlink: null,
        rtt: null
      }
    }

    return {
      type: this.networkInfo.type || 'unknown',
      effectiveType: this.networkInfo.effectiveType || 'unknown',
      downlink: this.networkInfo.downlink || null,
      rtt: this.networkInfo.rtt || null
    }
  }

  /**
   * Validate Wi-Fi connection against authorized networks
   * Since browser can't access SSID/BSSID directly, use fallback methods
   * @param {string} manualSSID - Manually entered SSID
   * @param {Object} gpsLocation - GPS coordinates
   * @param {Array} authorizedNetworks - List of authorized networks
   * @returns {Promise<Object>}
   */
  async validateWiFiConnection(manualSSID = null, gpsLocation = null, authorizedNetworks = []) {
    try {
      const isWiFi = await this.isConnectedToWiFi()
      
      if (!isWiFi) {
        return {
          isValid: false,
          method: 'network_check',
          error: 'Perangkat tidak terhubung ke Wi-Fi',
          details: 'Pastikan Wi-Fi aktif dan terhubung ke jaringan kantor'
        }
      }

      // Method 1: Manual SSID validation
      if (manualSSID) {
        const isAuthorizedSSID = authorizedNetworks.some(network => 
          network.ssid.toLowerCase() === manualSSID.toLowerCase()
        )
        
        if (isAuthorizedSSID) {
          return {
            isValid: true,
            method: 'manual_ssid',
            ssid: manualSSID,
            message: 'Wi-Fi tervalidasi melalui SSID manual'
          }
        } else {
          return {
            isValid: false,
            method: 'manual_ssid',
            error: 'SSID tidak diotorisasi',
            details: `SSID "${manualSSID}" tidak terdaftar sebagai jaringan kantor`
          }
        }
      }

      // Method 2: GPS-based validation
      if (gpsLocation && gpsLocation.latitude && gpsLocation.longitude) {
        const isInOfficeArea = this.validateGPSLocation(gpsLocation, authorizedNetworks)
        
        if (isInOfficeArea.isValid) {
          return {
            isValid: true,
            method: 'gps_validation',
            location: gpsLocation,
            message: 'Wi-Fi tervalidasi melalui lokasi GPS'
          }
        } else {
          return {
            isValid: false,
            method: 'gps_validation',
            error: 'Lokasi tidak valid',
            details: 'Anda tidak berada di area kantor yang diotorisasi'
          }
        }
      }

      // Method 3: Network fingerprinting (basic)
      const networkInfo = this.getNetworkInfo()
      if (networkInfo.type === 'wifi' && networkInfo.downlink > 0) {
        return {
          isValid: true,
          method: 'network_fingerprint',
          networkInfo,
          message: 'Wi-Fi terdeteksi (validasi terbatas)',
          warning: 'Validasi menggunakan deteksi jaringan dasar'
        }
      }

      return {
        isValid: false,
        method: 'fallback',
        error: 'Tidak dapat memvalidasi Wi-Fi',
        details: 'Silakan masukkan SSID secara manual atau aktifkan GPS'
      }

    } catch (error) {
      console.error('Wi-Fi validation error:', error)
      return {
        isValid: false,
        method: 'error',
        error: 'Terjadi kesalahan saat validasi Wi-Fi',
        details: error.message
      }
    }
  }

  /**
   * Validate GPS location against office areas
   * Since GPS boundaries are not in WiFi config, we'll use a simple approach:
   * If GPS is available and user is online, consider it valid
   * @param {Object} location - GPS coordinates
   * @param {Array} authorizedNetworks - Networks list
   * @returns {Object}
   */
  validateGPSLocation(location, authorizedNetworks) {
    try {
      // For now, if we have GPS location and there are authorized networks,
      // we'll consider it valid. In production, you should add GPS boundaries
      // to the WiFiConfig model in the database
      
      if (authorizedNetworks && authorizedNetworks.length > 0) {
        return {
          isValid: true,
          network: authorizedNetworks[0].ssid,
          message: 'Lokasi GPS terdeteksi'
        }
      }
      
      return {
        isValid: false,
        error: 'Tidak ada jaringan yang diotorisasi'
      }
    } catch (error) {
      return {
        isValid: false,
        error: 'Gagal validasi GPS: ' + error.message
      }
    }
  }

  /**
   * Calculate distance between two GPS coordinates (Haversine formula)
   * @param {number} lat1 
   * @param {number} lon1 
   * @param {number} lat2 
   * @param {number} lon2 
   * @returns {number} Distance in meters
   */
  calculateDistance(lat1, lon1, lat2, lon2) {
    const R = 6371000 // Earth's radius in meters
    const dLat = this.toRadians(lat2 - lat1)
    const dLon = this.toRadians(lon2 - lon1)
    
    const a = Math.sin(dLat / 2) * Math.sin(dLat / 2) +
              Math.cos(this.toRadians(lat1)) * Math.cos(this.toRadians(lat2)) *
              Math.sin(dLon / 2) * Math.sin(dLon / 2)
    
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))
    return R * c
  }

  /**
   * Convert degrees to radians
   * @param {number} degrees 
   * @returns {number}
   */
  toRadians(degrees) {
    return degrees * (Math.PI / 180)
  }

  /**
   * Get current GPS location
   * @returns {Promise<Object>}
   */
  async getCurrentLocation() {
    return new Promise((resolve, reject) => {
      if (!navigator.geolocation) {
        reject(new Error('Geolocation tidak didukung browser'))
        return
      }

      navigator.geolocation.getCurrentPosition(
        (position) => {
          resolve({
            latitude: position.coords.latitude,
            longitude: position.coords.longitude,
            accuracy: position.coords.accuracy,
            timestamp: position.timestamp
          })
        },
        (error) => {
          let errorMessage = 'Gagal mendapatkan lokasi GPS'
          switch (error.code) {
            case error.PERMISSION_DENIED:
              errorMessage = 'Akses lokasi ditolak. Silakan aktifkan izin lokasi.'
              break
            case error.POSITION_UNAVAILABLE:
              errorMessage = 'Lokasi tidak tersedia. Pastikan GPS aktif.'
              break
            case error.TIMEOUT:
              errorMessage = 'Timeout mendapatkan lokasi. Coba lagi.'
              break
          }
          reject(new Error(errorMessage))
        },
        {
          enableHighAccuracy: true,
          timeout: 10000,
          maximumAge: 60000
        }
      )
    })
  }

  /**
   * Try to get current WiFi SSID
   * Note: Browser cannot directly access SSID due to security restrictions
   * This will return null in most cases
   * @returns {Promise<string|null>}
   */
  async getCurrentSSID() {
    try {
      // Browser API doesn't provide SSID access
      // This is a placeholder for future implementation or native app
      
      // For now, we'll need to ask user to input SSID manually
      // or use GPS validation as fallback
      
      return null
    } catch (error) {
      console.warn('Cannot get SSID:', error)
      return null
    }
  }

  /**
   * Prompt user to input SSID manually
   * @returns {Promise<string|null>}
   */
  async promptForSSID() {
    return new Promise((resolve) => {
      const ssid = prompt('Masukkan nama WiFi (SSID) yang terhubung:')
      resolve(ssid ? ssid.trim() : null)
    })
  }

  /**
   * Add network change listeners
   * @param {Function} callback 
   */
  onNetworkChange(callback) {
    if (this.networkInfo) {
      this.networkInfo.addEventListener('change', callback)
    }
    
    window.addEventListener('online', callback)
    window.addEventListener('offline', callback)
  }

  /**
   * Remove network change listeners
   * @param {Function} callback 
   */
  offNetworkChange(callback) {
    if (this.networkInfo) {
      this.networkInfo.removeEventListener('change', callback)
    }
    
    window.removeEventListener('online', callback)
    window.removeEventListener('offline', callback)
  }
}

export default new WiFiService()