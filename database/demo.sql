-- 第一部分：SaaS 平台与租户管理 (SaaS Level)

-- 1. SaaS 平台管理用户表 (商家账号)
CREATE TABLE `platform_users`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `email`         varchar(255) NOT NULL UNIQUE COMMENT '登录邮箱',
    `phone`         varchar(20)  DEFAULT NULL COMMENT '联系电话',
    `password_hash` varchar(255) NOT NULL COMMENT '密码哈希',
    `real_name`     varchar(100) DEFAULT NULL COMMENT '真实姓名',
    `is_active`     tinyint(1) DEFAULT '1' COMMENT '是否激活',
    `created_at`    datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    `updated_at`    datetime(3) DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3),
    `deleted_at`    datetime(3) DEFAULT NULL,
    INDEX           `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='SaaS平台商家账号表';

-- 2. 订阅套餐表
CREATE TABLE `subscription_plans`
(
    `id`             bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name`           varchar(50)    NOT NULL COMMENT '套餐名称',
    `price_monthly`  decimal(12, 2) NOT NULL COMMENT '月付价格',
    `feature_limits` json DEFAULT NULL COMMENT '功能限制配置JSON',
    `created_at`     datetime(3) DEFAULT CURRENT_TIMESTAMP (3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='订阅套餐表';

-- 3. 组织/公司表 (新的计费主体)
CREATE TABLE `organizations`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name`       varchar(255) NOT NULL COMMENT '公司或团体名称',
    `type`       enum('individual', 'company') DEFAULT 'individual' COMMENT '主体类型',
    `tax_id`     varchar(100) DEFAULT NULL COMMENT '税号/企业统一代码',
    `owner_id`   bigint(20) unsigned NOT NULL COMMENT '组织创建者ID',
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    INDEX        `idx_owner` (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='组织主体表';

-- 4. 组织成员关系表 (支持多人协作)
CREATE TABLE `organization_members`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `org_id`     bigint(20) unsigned NOT NULL,
    `user_id`    bigint(20) unsigned NOT NULL,
    `role`       varchar(20) DEFAULT 'admin' COMMENT 'owner, admin, staff',
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    UNIQUE KEY `uk_org_user` (`org_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='组织成员关系表';

-- 5. 店铺租户表 (修正版)
CREATE TABLE `shops`
(
    `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `org_id`          bigint(20) unsigned NOT NULL COMMENT '所属组织ID',
    `plan_id`         bigint(20) unsigned NOT NULL COMMENT '订阅套餐ID',
    `name`            varchar(255) NOT NULL,
    `status`          varchar(20) DEFAULT 'active',
    `plan_expired_at` datetime(3) DEFAULT NULL,
    `config_settings` json        DEFAULT NULL,
    `created_at`      datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    `updated_at`      datetime(3) DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3),
    INDEX             `idx_org` (`org_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='店铺租户表';

-- 6. 店铺域名绑定表 (核心拓展)
CREATE TABLE `shop_domains`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `shop_id`    bigint(20) unsigned NOT NULL,
    `domain`     varchar(255) NOT NULL UNIQUE COMMENT '如 shop.abc.com 或 www.brand.com',
    `type`       enum('subdomain', 'custom') DEFAULT 'custom' COMMENT '二级域名或自定义独立域名',
    `is_primary` tinyint(1) DEFAULT '0' COMMENT '是否为主域名(唯一)',
    `ssl_status` varchar(20) DEFAULT 'pending' COMMENT 'SSL证书状态',
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    INDEX        `idx_shop` (`shop_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='店铺域名绑定表';


-- 第二部分：核心电商业务模块 (Core Commerce)
-- 7. 商品主表
CREATE TABLE `products`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `shop_id`    bigint(20) unsigned NOT NULL COMMENT '租户ID',
    `title`      varchar(255) NOT NULL COMMENT '标题',
    `body_html`  text COMMENT '详情描述',
    `status`     varchar(20) DEFAULT 'draft' COMMENT 'active, draft, archived',
    `metafields` json        DEFAULT NULL COMMENT '自定义属性/规格定义',
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3),
    INDEX        `idx_shop_status` (`shop_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='商品表';

-- 8. 商品规格与库存表 (SKU Level)
CREATE TABLE `product_variants`
(
    `id`                 bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `shop_id`            bigint(20) unsigned NOT NULL,
    `product_id`         bigint(20) unsigned NOT NULL,
    `sku`                varchar(100)   DEFAULT NULL,
    `price`              decimal(12, 2) NOT NULL,
    `compare_at_price`   decimal(12, 2) DEFAULT NULL COMMENT '原价/划线价',
    `inventory_quantity` int(11) DEFAULT '0',
    `option_values`      json           DEFAULT NULL COMMENT '规格选项值: {"color": "Red", "size": "XL"}',
    `created_at`         datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    INDEX                `idx_product_id` (`product_id`),
    INDEX                `idx_shop_sku` (`shop_id`, `sku`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='商品规格表';

-- =============================================================================
-- 第三部分：买家、购物车与弃单 (Buyer & Retention)
-- =============================================================================

-- 9. 店铺买家表
CREATE TABLE `customers`
(
    `id`                bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `shop_id`           bigint(20) unsigned NOT NULL,
    `email`             varchar(255) NOT NULL,
    `password_hash`     varchar(255)   DEFAULT NULL COMMENT '买家登录密码',
    `first_name`        varchar(100)   DEFAULT NULL,
    `last_name`         varchar(100)   DEFAULT NULL,
    `total_spent`       decimal(12, 2) DEFAULT '0.00',
    `accepts_marketing` tinyint(1) DEFAULT '0',
    `created_at`        datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    UNIQUE KEY `uk_shop_email` (`shop_id`, `email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='买家客户表';

-- 10. 购物车与弃单追踪
CREATE TABLE `carts`
(
    `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `shop_id`      bigint(20) unsigned NOT NULL,
    `token`        varchar(100) NOT NULL UNIQUE COMMENT '前端Cookie标识',
    `customer_id`  bigint(20) unsigned DEFAULT NULL,
    `items`        json DEFAULT NULL COMMENT '购物车快照',
    `is_abandoned` tinyint(1) DEFAULT '0' COMMENT '是否为弃单',
    `created_at`   datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    `updated_at`   datetime(3) DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='购物车与弃单表';



--  第三部分：物流、支付与营销 (Logistics, Payment & Marketing)

-- 11. 物流配置表
CREATE TABLE `shipping_rates`
(
    `id`                 bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `shop_id`            bigint(20) unsigned NOT NULL,
    `name`               varchar(100)   NOT NULL COMMENT '物流名称(如: 标准快递)',
    `price`              decimal(12, 2) NOT NULL COMMENT '运费',
    `min_order_subtotal` decimal(12, 2) DEFAULT NULL COMMENT '满多少包邮/可用',
    `countries`          json           DEFAULT NULL COMMENT '适用国家列表',
    `created_at`         datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    `updated_at`         datetime(3) DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3),
    `deleted_at`         datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY                  `idx_shop_id` (`shop_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='物流费率配置表';

-- 12. 支付项配置表
CREATE TABLE `payment_providers`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `shop_id`       bigint(20) unsigned NOT NULL,
    `provider_type` varchar(50) NOT NULL COMMENT '支付类型(Stripe, PayPal, Manual)',
    `config_data`   json DEFAULT NULL COMMENT '密钥、账号等加密配置',
    `is_enabled`    tinyint(1) DEFAULT '0',
    `created_at`    datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    `updated_at`    datetime(3) DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3),
    PRIMARY KEY (`id`) USING BTREE,
    KEY             `idx_shop_id` (`shop_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='支付方式配置表';

-- 13. 折扣/优惠券表
CREATE TABLE `discount_codes`
(
    `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `shop_id`         bigint(20) unsigned NOT NULL,
    `code`            varchar(50)    NOT NULL COMMENT '优惠码',
    `type`            varchar(20)    NOT NULL COMMENT 'percentage, fixed_amount, free_shipping',
    `value`           decimal(12, 2) NOT NULL COMMENT '折扣值',
    `min_requirement` decimal(12, 2) DEFAULT NULL COMMENT '最低消费金额',
    `starts_at`       datetime(3) DEFAULT NULL COMMENT '开始时间',
    `ends_at`         datetime(3) DEFAULT NULL COMMENT '结束时间',
    `usage_limit`     int(11) DEFAULT NULL COMMENT '总使用次数限制',
    `usage_count`     int(11) DEFAULT '0' COMMENT '当前已使用次数',
    `created_at`      datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    `updated_at`      datetime(3) DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3),
    `deleted_at`      datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY               `idx_shop_code` (`shop_id`, `code`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='折扣/优惠券表';

-- 14. 支付流水表
CREATE TABLE `payment_transactions`
(
    `id`               bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `shop_id`          bigint(20) unsigned NOT NULL,
    `order_id`         bigint(20) unsigned NOT NULL,
    `transaction_type` varchar(20)    NOT NULL COMMENT 'sale, refund',
    `gateway`          varchar(50)    NOT NULL COMMENT 'stripe, paypal',
    `gateway_ref`      varchar(255) DEFAULT NULL COMMENT '网关流水号',
    `amount`           decimal(12, 2) NOT NULL,
    `status`           varchar(20)  DEFAULT 'pending' COMMENT 'pending, success, failed',
    `raw_response`     json         DEFAULT NULL COMMENT '原始回执',
    `created_at`       datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    INDEX              `idx_order_gateway` (`order_id`, `gateway_ref`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='支付流水表';



--- 第四部分：内容与展示模块 (CMS & Themes)

-- 15. 博客推文表
CREATE TABLE `blog_posts`
(
    `id`             bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `shop_id`        bigint(20) unsigned NOT NULL,
    `title`          varchar(255) NOT NULL,
    `author`         varchar(100) DEFAULT NULL,
    `content_html`   longtext,
    `summary`        text COMMENT '摘要',
    `featured_image` varchar(255) DEFAULT NULL COMMENT '封面图',
    `status`         varchar(20)  DEFAULT 'published' COMMENT 'published, draft',
    `published_at`   datetime(3) DEFAULT NULL,
    `created_at`     datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    `updated_at`     datetime(3) DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3),
    `deleted_at`     datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY              `idx_shop_status` (`shop_id`, `status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='博客推文表';

-- 16. 模板/主题表
CREATE TABLE `themes`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `shop_id`     bigint(20) unsigned NOT NULL,
    `name`        varchar(100) DEFAULT NULL COMMENT '模板名',
    `source_path` varchar(255) DEFAULT NULL COMMENT '模板文件S3/本地路径',
    `config_data` json         DEFAULT NULL COMMENT '模板颜色、布局等JSON配置',
    `is_active`   tinyint(1) DEFAULT '0' COMMENT '当前激活模板',
    `created_at`  datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    `updated_at`  datetime(3) DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3),
    `deleted_at`  datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY           `idx_shop_id` (`shop_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='模板/主题表';


--  第五部分：订单与流水 (Transactions)

-- 17. 订单表
CREATE TABLE `orders`
(
    `id`                 bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `shop_id`            bigint(20) unsigned NOT NULL,
    `order_number`       varchar(50)    NOT NULL COMMENT '对外展示订单号 #1001',

    -- 1. 客户身份 (扩展)
    `customer_id`        bigint(20) unsigned DEFAULT NULL COMMENT '关联注册买家ID(游客为空)',
    `customer_email`     varchar(255)   NOT NULL COMMENT '下单邮箱(快照)',
    `customer_phone`     varchar(50)    DEFAULT NULL COMMENT '下单电话(快照)',

    -- 2. 财务信息
    `currency`           varchar(10)    NOT NULL COMMENT 'USD, CNY',
    `total_price`        decimal(12, 2) NOT NULL COMMENT '应付总额(含税运)',
    `subtotal_price`     decimal(12, 2) NOT NULL COMMENT '商品小计(不含税运)',
    `total_tax`          decimal(12, 2) DEFAULT '0.00' COMMENT '税费总额',
    `total_discounts`    decimal(12, 2) DEFAULT '0.00' COMMENT '折扣总额',
    `shipping_price`     decimal(12, 2) DEFAULT '0.00' COMMENT '运费',

    -- 3. 状态管理
    `financial_status`   varchar(20)    DEFAULT 'pending' COMMENT 'paid, refunded, voided',
    `fulfillment_status` varchar(20)    DEFAULT 'unfulfilled' COMMENT 'shipped, partial',
    `cancel_reason`      varchar(50)    DEFAULT NULL COMMENT 'customer, fraud, inventory, other',

    -- 4. 地址快照 (核心补充：不要关联ID，直接存JSON)
    -- 格式: {"first_name": "John", "address1": "123 Main St", "city": "NY", "country": "US", "zip": "10001"}
    `shipping_address`   json           DEFAULT NULL COMMENT '收货地址快照',
    `billing_address`    json           DEFAULT NULL COMMENT '账单地址快照',

    -- 5. 审计与环境 (核心补充：风控必备)
    `note`               text COMMENT '客户留言',
    `tags`               varchar(255)   DEFAULT NULL COMMENT '商家标签: VIP, 疑似欺诈',
    `client_ip`          varchar(45)    DEFAULT NULL COMMENT '下单IP',
    `user_agent`         text COMMENT '用户设备信息',
    `landing_site`       varchar(255)   DEFAULT NULL COMMENT '落地页(广告来源追踪)',

    -- 6. 时间记录
    `processed_at`       datetime(3) DEFAULT NULL COMMENT '支付成功/订单确认时间',
    `created_at`         datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    `updated_at`         datetime(3) DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3),
    `deleted_at`         datetime(3) DEFAULT NULL, -- 逻辑删除

    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_shop_order` (`shop_id`, `order_number`) USING BTREE,
    KEY                  `idx_customer` (`shop_id`, `customer_id`),
    KEY                  `idx_created` (`shop_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='订单主表';

-- 18. 订单详情快照
CREATE TABLE `order_items`
(
    `id`                   bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `shop_id`              bigint(20) unsigned NOT NULL,
    `order_id`             bigint(20) unsigned NOT NULL,
    `product_id`           bigint(20) unsigned DEFAULT NULL COMMENT '关联商品ID(商品删了也能留空)',
    `variant_id`           bigint(20) unsigned DEFAULT NULL COMMENT '关联规格ID',

    -- 1. 基础信息快照 (冗余字段，方便快速查询和显示)
    `name`                 varchar(255)   NOT NULL COMMENT '商品名称快照',
    `sku`                  varchar(100)   DEFAULT NULL COMMENT 'SKU快照(方便仓库扫码)',
    `image_url`            varchar(512)   DEFAULT NULL COMMENT '商品主图快照',

    -- 2. 数量与价格
    `quantity`             int(11) NOT NULL COMMENT '购买数量',
    `fulfillable_quantity` int(11) NOT NULL COMMENT '待发货数量(处理部分发货)',
    `price`                decimal(12, 2) NOT NULL COMMENT '单价(快照)',
    `total_discount`       decimal(12, 2) DEFAULT '0.00' COMMENT '分摊到该商品的折扣金额',

    -- 3. 深度快照 (核心补充)
    -- 格式: {"options": [{"name": "Size", "value": "M"}], "weight": 0.5, "compare_at_price": 20.00}
    `variant_snapshot`     json           DEFAULT NULL COMMENT '完整规格数据快照',

    -- 4. 定制化属性 (核心补充：SaaS 插件通过这里存数据)
    -- 格式: [{"name": "刻字", "value": "Happy Birthday"}, {"name": "上传图片", "value": "s3://..."}]
    `properties`           json           DEFAULT NULL COMMENT '用户自定义属性(Custom Fields)',

    `created_at`           datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    PRIMARY KEY (`id`) USING BTREE,
    KEY                    `idx_order` (`order_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='订单明细表';

-- 19. 库存流水
CREATE TABLE `inventory_histories`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `shop_id`       bigint(20) unsigned NOT NULL,
    `variant_id`    bigint(20) unsigned NOT NULL,
    `change_amount` int(11) NOT NULL,
    `reason`        varchar(100) DEFAULT NULL,
    `reference_id`  varchar(100) DEFAULT NULL COMMENT '订单ID或退款单ID',
    `created_at`    datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    PRIMARY KEY (`id`) USING BTREE,
    KEY             `idx_variant_id` (`variant_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='库存变更记录表';


-- 20. 店铺支持的语言表
CREATE TABLE `shop_languages`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `shop_id`    bigint(20) unsigned NOT NULL COMMENT '租户ID',
    `locale`     varchar(10) NOT NULL COMMENT '语言代码(如: zh-CN, en-US)',
    `name`       varchar(50) NOT NULL COMMENT '显示名称(如: 简体中文)',
    `is_default` tinyint(1) DEFAULT '0' COMMENT '是否为默认语言',
    `is_enabled` tinyint(1) DEFAULT '1' COMMENT '是否启用',
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3),
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_shop_locale` (`shop_id`, `locale`) USING BTREE,
    KEY          `idx_shop_id` (`shop_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='店铺语言配置表';

-- 21. 店铺支持的货币表
CREATE TABLE `shop_currencies`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `shop_id`       bigint(20) unsigned NOT NULL COMMENT '租户ID',
    `currency_code` varchar(10) NOT NULL COMMENT '货币代码(如: CNY, USD, EUR)',
    `symbol`        varchar(10)    DEFAULT NULL COMMENT '货币符号(如: ¥, $)',
    `exchange_rate` decimal(18, 6) DEFAULT '1.000000' COMMENT '相对于主货币的汇率',
    `is_default`    tinyint(1) DEFAULT '0' COMMENT '是否为默认结算货币',
    `is_enabled`    tinyint(1) DEFAULT '1' COMMENT '是否启用',
    `created_at`    datetime(3) DEFAULT CURRENT_TIMESTAMP (3),
    `updated_at`    datetime(3) DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3),
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_shop_currency` (`shop_id`, `currency_code`) USING BTREE,
    KEY             `idx_shop_id` (`shop_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='店铺货币配置表';