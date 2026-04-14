package com.legalsys.dto;

import lombok.Data;
import java.util.List;

@Data
public class PageResponse<T> {
    private List<T> items;
    private long total;
    private int page;
    private int pageSize;

    public static <T> PageResponse<T> of(List<T> items, long total, int page, int pageSize) {
        PageResponse<T> response = new PageResponse<>();
        response.setItems(items);
        response.setTotal(total);
        response.setPage(page);
        response.setPageSize(pageSize);
        return response;
    }
}
