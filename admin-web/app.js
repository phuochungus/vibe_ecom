const products = [
  { name: 'Gậy Driver TaylorMade Stealth', category: 'Gậy golf', price: 12900000, stock: 12 },
  { name: 'Bóng Titleist Pro V1 (12 bóng)', category: 'Bóng golf', price: 1250000, stock: 34 },
  { name: 'Găng tay FootJoy WeatherSof', category: 'Phụ kiện', price: 520000, stock: 50 },
];

const orders = [
  { id: 'ORD-1001', customer: 'Nguyễn Văn A', value: 14500000, status: 'Đang xử lý' },
  { id: 'ORD-1002', customer: 'Trần Thị B', value: 2190000, status: 'Hoàn tất' },
  { id: 'ORD-1003', customer: 'Lê Minh C', value: 8700000, status: 'Mới' },
];

const customers = ['Phạm Quốc Dũng', 'Ngô Hải Yến', 'Đào Thanh Tuấn', 'Trần Minh Anh'];

const adminUsers = [
  { name: 'Vũ Đức Hoàng', email: 'hoang.vu@golfstore.vn', role: 'Quản trị viên', active: true },
  { name: 'Lâm Hải My', email: 'my.lam@golfstore.vn', role: 'Biên tập viên', active: true },
  { name: 'Phan Gia Khánh', email: 'khanh.phan@golfstore.vn', role: 'Hỗ trợ', active: false },
];

const sections = document.querySelectorAll('.panel');
const navButtons = document.querySelectorAll('.nav-btn');
const titleEl = document.getElementById('section-title');
const productTable = document.getElementById('product-table');
const orderList = document.getElementById('order-list');
const customerList = document.getElementById('customer-list');
const adminTable = document.getElementById('admin-table');

const monthlyRevenue = document.getElementById('monthly-revenue');
const newOrders = document.getElementById('new-orders');
const activeAdmins = document.getElementById('active-admins');

function formatCurrency(value) {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    maximumFractionDigits: 0,
  }).format(value);
}

function renderProducts(list = products) {
  productTable.innerHTML = list
    .map(
      (product) => `
      <tr>
        <td>${product.name}</td>
        <td>${product.category}</td>
        <td>${formatCurrency(product.price)}</td>
        <td>${product.stock}</td>
      </tr>
    `,
    )
    .join('');
}

function renderOrders(list = orders) {
  orderList.innerHTML = list
    .map((order) => `<li><strong>${order.id}</strong> - ${order.customer} - ${formatCurrency(order.value)} (${order.status})</li>`)
    .join('');
}

function renderCustomers(list = customers) {
  customerList.innerHTML = list.map((name) => `<li>${name}</li>`).join('');
}

function renderAdminUsers(list = adminUsers) {
  adminTable.innerHTML = list
    .map(
      (admin, index) => `
      <tr>
        <td>${admin.name}</td>
        <td>${admin.email}</td>
        <td>${admin.role}</td>
        <td><span class="badge ${admin.active ? 'active' : 'inactive'}">${admin.active ? 'Đang hoạt động' : 'Tạm khóa'}</span></td>
        <td><button class="action-btn" data-toggle-admin="${index}">${admin.active ? 'Khóa' : 'Mở khóa'}</button></td>
      </tr>
    `,
    )
    .join('');
}

function updateDashboard() {
  const revenue = orders.reduce((sum, order) => sum + order.value, 0);
  monthlyRevenue.textContent = formatCurrency(revenue);
  newOrders.textContent = orders.filter((order) => order.status === 'Mới').length;
  activeAdmins.textContent = adminUsers.filter((admin) => admin.active).length;
}

navButtons.forEach((btn) => {
  btn.addEventListener('click', () => {
    const target = btn.dataset.section;

    navButtons.forEach((item) => item.classList.remove('active'));
    btn.classList.add('active');

    sections.forEach((section) => {
      section.classList.toggle('active', section.id === target);
    });

    titleEl.textContent = btn.textContent;
  });
});

document.getElementById('product-form').addEventListener('submit', (event) => {
  event.preventDefault();

  const name = document.getElementById('product-name').value.trim();
  const price = Number(document.getElementById('product-price').value);

  if (!name || Number.isNaN(price) || price < 0) {
    return;
  }

  products.unshift({ name, price, stock: 0, category: 'Khác' });
  renderProducts();
  event.target.reset();
});

document.getElementById('admin-form').addEventListener('submit', (event) => {
  event.preventDefault();

  const name = document.getElementById('admin-name').value.trim();
  const email = document.getElementById('admin-email').value.trim();
  const role = document.getElementById('admin-role').value;

  if (!name || !email || !role) {
    return;
  }

  adminUsers.unshift({ name, email, role, active: true });
  renderAdminUsers();
  updateDashboard();
  event.target.reset();
});

adminTable.addEventListener('click', (event) => {
  const { toggleAdmin } = event.target.dataset;
  if (toggleAdmin === undefined) {
    return;
  }

  const index = Number(toggleAdmin);
  if (Number.isNaN(index) || !adminUsers[index]) {
    return;
  }

  adminUsers[index].active = !adminUsers[index].active;
  renderAdminUsers();
  updateDashboard();
});

document.getElementById('search').addEventListener('input', (event) => {
  const keyword = event.target.value.trim().toLowerCase();

  if (!keyword) {
    renderProducts();
    renderOrders();
    renderCustomers();
    renderAdminUsers();
    return;
  }

  const filteredProducts = products.filter(
    (product) =>
      product.name.toLowerCase().includes(keyword) ||
      product.category.toLowerCase().includes(keyword),
  );

  const filteredOrders = orders.filter(
    (order) =>
      order.id.toLowerCase().includes(keyword) ||
      order.customer.toLowerCase().includes(keyword) ||
      order.status.toLowerCase().includes(keyword),
  );

  const filteredCustomers = customers.filter((name) => name.toLowerCase().includes(keyword));

  const filteredAdmins = adminUsers.filter(
    (admin) =>
      admin.name.toLowerCase().includes(keyword) ||
      admin.email.toLowerCase().includes(keyword) ||
      admin.role.toLowerCase().includes(keyword),
  );

  renderProducts(filteredProducts);
  renderOrders(filteredOrders);
  renderCustomers(filteredCustomers);
  renderAdminUsers(filteredAdmins);
});

renderProducts();
renderOrders();
renderCustomers();
renderAdminUsers();
updateDashboard();
