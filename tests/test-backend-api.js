const axios = require('axios');

async function testBackendAPI() {
  const API_BASE_URL = 'http://localhost:8080/api/v1';
  
  console.log('\n=== Testing Backend API ===\n');

  try {
    // Step 1: Login to get token
    console.log('1. Logging in...');
    const loginResponse = await axios.post(`${API_BASE_URL}/auth/login`, {
      identifier: 'kepala.sppg@sppg.com',
      password: 'password123'
    });

    if (!loginResponse.data.success) {
      console.log('✗ Login failed:', loginResponse.data);
      return;
    }

    const token = loginResponse.data.token;
    console.log('✓ Login successful');
    console.log(`  Token: ${token.substring(0, 20)}...`);

    // Step 2: Create a supplier
    console.log('\n2. Creating supplier...');
    const supplierData = {
      name: `API Test Supplier ${Date.now()}`,
      product_category: 'Sayuran',
      contact_person: 'API Test Contact',
      phone_number: '081234567890',
      email: `apitest${Date.now()}@test.com`,
      address: 'Jl. API Test No. 123'
    };

    console.log('  Data:', JSON.stringify(supplierData, null, 2));

    const createResponse = await axios.post(
      `${API_BASE_URL}/suppliers`,
      supplierData,
      {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      }
    );

    console.log('\n  Response status:', createResponse.status);
    console.log('  Response data:', JSON.stringify(createResponse.data, null, 2));

    if (createResponse.data.success) {
      console.log('✓ Supplier created successfully');
      const supplierId = createResponse.data.supplier?.id;
      console.log(`  Supplier ID: ${supplierId}`);

      // Step 3: Fetch all suppliers to verify
      console.log('\n3. Fetching all suppliers...');
      const listResponse = await axios.get(`${API_BASE_URL}/suppliers`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      console.log('  Total suppliers:', listResponse.data.suppliers?.length || 0);
      
      // Check if our supplier is in the list
      const ourSupplier = listResponse.data.suppliers?.find(s => s.id === supplierId);
      if (ourSupplier) {
        console.log('✓ Supplier found in list');
        console.log('  Supplier:', JSON.stringify(ourSupplier, null, 2));
      } else {
        console.log('✗ Supplier NOT found in list');
        console.log('  All suppliers:', JSON.stringify(listResponse.data.suppliers, null, 2));
      }

      // Step 4: Fetch the specific supplier
      if (supplierId) {
        console.log('\n4. Fetching specific supplier...');
        const getResponse = await axios.get(`${API_BASE_URL}/suppliers/${supplierId}`, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        console.log('  Response:', JSON.stringify(getResponse.data, null, 2));
      }

    } else {
      console.log('✗ Supplier creation failed');
      console.log('  Error:', createResponse.data);
    }

  } catch (error) {
    console.error('\n✗ Error:', error.message);
    if (error.response) {
      console.error('  Status:', error.response.status);
      console.error('  Data:', JSON.stringify(error.response.data, null, 2));
    }
  }

  console.log('\n=== Test Complete ===\n');
}

testBackendAPI();
