document.getElementById('recipeForm').addEventListener('submit', function(event) {
    event.preventDefault(); // Предотвращаем отправку формы по умолчанию

    const formData = new FormData(this); // Создаем объект FormData из формы

    fetch('/create-recipe', {
        method: 'POST',
        body: formData // Отправляем данные формы на сервер методом POST
    })
        .then(response => response.json()) // Получаем ответ от сервера в формате JSON
        .then(data => {
            console.log('Recipe created successfully:', data);
            // Дальнейшие действия, например, обновление интерфейса или переход на другую страницу
        })
        .catch(error => {
            console.error('Error creating recipe:', error);
            // Обработка ошибок при создании рецепта
        });
});

function addIngredient() {
    const ingredientFields = document.getElementById('ingredientFields');
    const newIngredient = document.createElement('div');
    newIngredient.className = 'ingredient';
    newIngredient.innerHTML = `
        <input type="text" name="ingredients[]" placeholder="Название ингредиента" required>
        <input type="number" name="quantities[]" placeholder="Количество" required>
        <select name="units[]">
            <option value="шт">шт</option>
            <option value="г">г</option>
            <option value="мл">мл</option>
        </select>
    `;
    ingredientFields.appendChild(newIngredient);
}
