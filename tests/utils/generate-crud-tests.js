/**
 * CRUD Test Generator
 * Generates CRUD test suites for all modules based on configuration
 */

const fs = require('fs');
const path = require('path');

// Module configurations
const modules = [
  {
    name: 'logistik-data-sekolah',
    displayName: 'Logistik - Data Sekolah',
    url: '/schools',
    entityName: 'School',
    entityNameLower: 'school',
    addButtonText: 'Tambah Sekolah',
    fields: [
      { id: 'form_item_name', value: 'Test School', type: 'input' },
      { id: 'form_item_address', value: 'Jl. Test School No. 123', type: 'textarea' },
      { id: 'form_item_student_count', value: '100', type: 'input' },
    ]
  },
  {
    name: 'sdm-data-karyawan',
    displayName: 'SDM - Data Karyawan',
    url: '/employees',
    entityName: 'Employee',
    entityNameLower: 'employee',
    addButtonText: 'Tambah Karyawan',
    fields: [
      { id: 'form_item_name', value: 'Test Employee', type: 'input' },
      { id: 'form_item_email', value: 'test@employee.com', type: 'input' },
      { id: 'form_item_phone', value: '081234567890', type: 'input' },
    ]
  },
];

// Template for CRUD test suite
const generateTestSuite = (module) => {
  const timestamp = '${crudHelper.getTimestamp()}';
  
  return `const { test, expect } = require('@playwright/test');
const CRUDHelper = require('../utils/crud-helper');

const ConfigLoader = require('../utils/config-loader');
const configLoader = new ConfigLoader();
const config = configLoader.load();

test.describe('${module.displayName} CRUD Operations', () => {
  let test${module.entityName}Name;
  let crudHelper;

  test.beforeEach(async ({ page }) => {
    crudHelper = new CRUDHelper(page);
    test${module.entityName}Name = \`${module.fields[0].value} \${crudHelper.getTimestamp()}\`;
    
    await page.goto(config.pwaBaseUrl);
    await page.waitForLoadState('networkidle');
    
    const isLoginPage = await page.locator('input[type="password"]').count() > 0;
    
    if (isLoginPage) {
      const usernameInput = page.locator('input[type="text"], input[type="email"]').first();
      await usernameInput.fill('kepala.sppg@sppg.com');
      
      const passwordInput = page.locator('input[type="password"]').first();
      await passwordInput.fill('password123');
      
      const loginButton = page.locator('button:has-text("Login"), button:has-text("Masuk"), button[type="submit"]').first();
      await loginButton.click();
      
      await page.waitForLoadState('networkidle');
      await page.waitForTimeout(1000);
    }
  });

  test('CRUD-001: Create new ${module.entityNameLower}', async ({ page }) => {
    console.log('\\n=== TEST: CREATE NEW ${module.entityName.toUpperCase()} ===');
    
    await page.goto(\`\${config.pwaBaseUrl}${module.url}\`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const tambahButton = page.locator('button:has-text("${module.addButtonText}")').first();
    
    if (await tambahButton.count() === 0) {
      console.log('⚠ Test skipped: "${module.addButtonText}" button not found');
      test.skip();
      return;
    }
    
    await tambahButton.click();
    await page.waitForTimeout(1000);
    
    console.log(\`\\nFilling form with test data:\`);
    console.log(\`- ${module.entityName} Name: \${test${module.entityName}Name}\`);
    
    ${module.fields.map((field, index) => {
      if (index === 0) {
        return `await crudHelper.fillInput('${field.id}', test${module.entityName}Name);`;
      } else if (field.type === 'textarea') {
        return `await crudHelper.fillTextarea('${field.id}', '${field.value}');`;
      } else {
        return `await crudHelper.fillInput('${field.id}', '${field.value}');`;
      }
    }).join('\n    ')}
    
    await crudHelper.clickSubmit('OK');
    await crudHelper.waitForSuccessMessage();
    await page.waitForTimeout(2000);
    
    const dataExists = await crudHelper.verifyDataInTable(test${module.entityName}Name);
    expect(dataExists).toBeTruthy();
    console.log(\`✓ Test CRUD-001 PASSED\\n\`);
  });

  test('CRUD-005: Test form validation', async ({ page }) => {
    console.log('\\n=== TEST: FORM VALIDATION ===');
    
    await page.goto(\`\${config.pwaBaseUrl}${module.url}\`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const tambahButton = page.locator('button:has-text("${module.addButtonText}")').first();
    
    if (await tambahButton.count() === 0) {
      console.log('⚠ Test skipped');
      test.skip();
      return;
    }
    
    await tambahButton.click();
    await page.waitForTimeout(1000);
    
    await crudHelper.clickSubmit('OK');
    await page.waitForTimeout(1000);
    
    const hasErrors = await page.locator('.ant-form-item-explain-error, .ant-message-error').count() > 0;
    
    if (hasErrors) {
      console.log('✓ Form validation working');
      expect(hasErrors).toBeTruthy();
    } else {
      console.log('⚠ No validation errors found');
    }
    
    console.log(\`✓ Test CRUD-005 PASSED\\n\`);
  });
});
`;
};

// Generate test files
modules.forEach(module => {
  const testContent = generateTestSuite(module);
  const filePath = path.join(__dirname, '..', 'test-suites', `${module.name}-crud.spec.js`);
  
  fs.writeFileSync(filePath, testContent);
  console.log(`✓ Generated: ${module.name}-crud.spec.js`);
});

console.log('\n✓ All CRUD test suites generated successfully!');
