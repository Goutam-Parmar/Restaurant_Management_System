BEGIN;

TRUNCATE TABLE
    payments,
  order_items,
  orders,
  menus,
  restaurants,
  addresses,
  user_roles,
  users
RESTART IDENTITY CASCADE;

COMMIT;
