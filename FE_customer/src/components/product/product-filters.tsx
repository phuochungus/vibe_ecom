import { Card, Slider, Select, InputNumber, Button, Space, Divider } from "antd";
import { useState } from "react";

interface ProductFiltersProps {
    onFilterChange: (filters: { minPrice?: number; maxPrice?: number; sort?: string }) => void;
}

export default function ProductFilters({ onFilterChange }: ProductFiltersProps) {
    const [priceRange, setPriceRange] = useState<[number, number]>([0, 50000000]);
    const [sort, setSort] = useState<string>();

    const handleApply = () => {
        onFilterChange({
            minPrice: priceRange[0],
            maxPrice: priceRange[1],
            sort,
        });
    };

    const handleReset = () => {
        setPriceRange([0, 50000000]);
        setSort(undefined);
        onFilterChange({});
    };

    return (
        <Card title="Bộ lọc" size="small">
            <Space orientation="vertical" style={{ width: "100%" }} size="large">
                {/* Price Range */}
                <div>
                    <div style={{ marginBottom: 8, fontWeight: 500 }}>Khoảng giá</div>
                    <Slider
                        range
                        min={0}
                        max={50000000}
                        step={100000}
                        value={priceRange}
                        onChange={(value) => setPriceRange(value as [number, number])}
                        tooltip={{
                            formatter: (value) => `${(value! / 1000000).toFixed(1)}M`,
                        }}
                    />
                    <Space>
                        <InputNumber
                            value={priceRange[0]}
                            onChange={(val) => setPriceRange([val || 0, priceRange[1]])}
                            formatter={(value) => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ",")}
                            parser={(value) => value!.replace(/,/g, "") as unknown as number}
                            style={{ width: 120 }}
                            placeholder="Từ"
                        />
                        <span>-</span>
                        <InputNumber
                            value={priceRange[1]}
                            onChange={(val) => setPriceRange([priceRange[0], val || 50000000])}
                            formatter={(value) => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ",")}
                            parser={(value) => value!.replace(/,/g, "") as unknown as number}
                            style={{ width: 120 }}
                            placeholder="Đến"
                        />
                    </Space>
                </div>

                <Divider style={{ margin: "8px 0" }} />

                {/* Sort */}
                <div>
                    <div style={{ marginBottom: 8, fontWeight: 500 }}>Sắp xếp</div>
                    <Select
                        value={sort}
                        onChange={setSort}
                        style={{ width: "100%" }}
                        placeholder="Chọn cách sắp xếp"
                        allowClear
                        options={[
                            { label: "Giá: Thấp đến cao", value: "price_asc" },
                            { label: "Giá: Cao đến thấp", value: "price_desc" },
                            { label: "Mới nhất", value: "newest" },
                            { label: "Cũ nhất", value: "oldest" },
                        ]}
                    />
                </div>

                {/* Actions */}
                <Space style={{ width: "100%", justifyContent: "space-between" }}>
                    <Button onClick={handleReset}>Xóa bộ lọc</Button>
                    <Button type="primary" onClick={handleApply}>
                        Áp dụng
                    </Button>
                </Space>
            </Space>
        </Card>
    );
}
