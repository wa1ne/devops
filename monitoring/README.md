## Инфа о дашборде
1. [Ссылка на дашборд](monitoring.shp-devops.run.place:30036/d/fegpeotkrsgzkf/devops-client-by-menshikov-a)
2. ![Скриншот дашборда](screenshot.png)
3. Используемые сигналы
    - requested_total
    - response_time_bucket с перцентилями P50, P75, P90, P99
    - image_request_total c фильтром need_image="image_requested"
    - requested_types_total с филтрами type="traffic1/2/3"