document.addEventListener('DOMContentLoaded', function() {
    loadInstructors();
    setupBookingFlow();
});

function loadInstructors() {
    fetch('/api/instructors')
        .then(function(res) { return res.json(); })
        .then(function(instructors) {
            var grid = document.getElementById('instructors-grid');
            if (!grid) return;
            if (instructors.length === 0) {
                grid.innerHTML = '<p class="text-gray-500 text-center col-span-full">Инструкторы скоро появятся</p>';
                return;
            }
            var placeholder = 'https://via.placeholder.com/400x400?text=Instructor';
            instructors.forEach(function(instructor) {
                var photoURL = instructor.Photo || placeholder;
                var div = document.createElement('div');
                div.className = 'bg-white rounded-lg shadow-lg overflow-hidden hover:shadow-xl transition';

                var img = document.createElement('img');
                img.src = photoURL;
                img.alt = instructor.Name;
                img.className = 'w-full h-64 object-cover';
                img.onerror = function() { this.src = placeholder; };
                div.appendChild(img);

                var body = document.createElement('div');
                body.className = 'p-6';
                body.innerHTML = '<h3 class="text-2xl font-bold text-gray-800 mb-2">' + instructor.Name + '</h3>' +
                    '<p class="text-gray-600 mb-3">' + (instructor.Phone || '') + '</p>' +
                    '<p class="text-gray-700 mb-3">' + (instructor.Description || 'Опытный инструктор SUP') + '</p>';

                if (instructor.WalkTypes && instructor.WalkTypes.length > 0) {
                    var wtHtml = '<div class="text-sm text-gray-600"><p class="font-semibold mb-1">Прогулки:</p>';
                    instructor.WalkTypes.forEach(function(wt) {
                        wtHtml += '<p>• ' + wt.Name + ' — ' + wt.Price + ' ₽, до ' + wt.MaxPeople + ' чел.</p>';
                    });
                    wtHtml += '</div>';
                    body.innerHTML += wtHtml;
                }

                div.appendChild(body);
                grid.appendChild(div);
            });
        });
}

var bookingState = {
    instructor: null,
    walkType: null,
    slot: null
};

function setupBookingFlow() {
    if (!document.getElementById('booking-instructors')) return;
    loadBookingInstructors();

    var form = document.getElementById('booking-form');
    if (!form) return;
    form.addEventListener('submit', submitBookingForm);
}

function loadBookingInstructors() {
    fetch('/api/instructors')
        .then(function(res) { return res.json(); })
        .then(function(instructors) {
            var container = document.getElementById('booking-instructors');
            if (instructors.length === 0) {
                container.innerHTML = '<p class="text-gray-500">Инструкторы пока недоступны</p>';
                return;
            }
            container.innerHTML = '';
            instructors.forEach(function(inst) {
                var card = document.createElement('button');
                card.type = 'button';
                card.className = 'text-left border rounded-lg p-4 hover:border-blue-500';
                card.innerHTML = '<p class="font-semibold">' + inst.Name + '</p><p class="text-sm text-gray-600">' + (inst.Description || '') + '</p>';
                card.onclick = function() { selectInstructor(inst); };
                container.appendChild(card);
            });
        });
}

function selectInstructor(inst) {
    bookingState.instructor = inst;
    bookingState.walkType = null;
    bookingState.slot = null;
    document.getElementById('walk-type-step').classList.remove('hidden');
    document.getElementById('slot-step').classList.add('hidden');
    document.getElementById('booking-form-container').classList.add('hidden');

    fetch('/api/instructors/' + inst.ID + '/walk-types')
        .then(function(res) { return res.json(); })
        .then(function(walkTypes) {
            var container = document.getElementById('walk-types-container');
            container.innerHTML = '';
            if (walkTypes.length === 0) {
                container.innerHTML = '<p class="text-gray-500">У инструктора пока нет типов прогулок</p>';
                return;
            }
            walkTypes.forEach(function(wt) {
                var card = document.createElement('button');
                card.type = 'button';
                card.className = 'text-left border rounded-lg p-4 hover:border-blue-500';
                card.innerHTML = '<p class="font-semibold">' + wt.Name + '</p>' +
                    '<p class="text-sm text-blue-700">' + wt.Price + ' ₽</p>' +
                    '<p class="text-sm text-gray-600">до ' + wt.MaxPeople + ' чел.</p>';
                card.onclick = function() { selectWalkType(wt); };
                container.appendChild(card);
            });
        });
}

function selectWalkType(walkType) {
    bookingState.walkType = walkType;
    bookingState.slot = null;
    document.getElementById('slot-step').classList.remove('hidden');
    document.getElementById('booking-form-container').classList.add('hidden');

    fetch('/api/slots?instructor_id=' + bookingState.instructor.ID + '&walk_type_id=' + walkType.ID)
        .then(function(res) { return res.json(); })
        .then(function(slots) {
            var container = document.getElementById('slots-container');
            container.innerHTML = '';
            if (slots.length === 0) {
                container.innerHTML = '<p class="text-gray-500">Нет доступных слотов для выбранной прогулки</p>';
                return;
            }

            var grouped = {};
            slots.forEach(function(slot) {
                if (slot.Status !== 'available') return;
                var date = new Date(slot.Date).toLocaleDateString('ru-RU');
                if (!grouped[date]) grouped[date] = [];
                grouped[date].push(slot);
            });

            Object.keys(grouped).sort().forEach(function(date) {
                var dateDiv = document.createElement('div');
                dateDiv.className = 'mb-4';
                dateDiv.innerHTML = '<h3 class="text-lg font-semibold mb-2">' + date + '</h3>';
                var grid = document.createElement('div');
                grid.className = 'grid grid-cols-1 md:grid-cols-2 gap-3';

                grouped[date].forEach(function(slot) {
                    var btn = document.createElement('button');
                    btn.type = 'button';
                    btn.className = 'border rounded-lg p-4 text-left hover:border-blue-500';
                    btn.innerHTML = '<p class="font-semibold">' + slot.StartTime.substring(0, 5) + ' - ' + slot.EndTime.substring(0, 5) + '</p>' +
                        '<p class="text-sm text-gray-600">' + slot.Price + ' ₽ • до ' + slot.MaxPeople + ' чел.</p>';
                    btn.onclick = function() { selectSlot(slot, date); };
                    grid.appendChild(btn);
                });

                dateDiv.appendChild(grid);
                container.appendChild(dateDiv);
            });
        });
}

function selectSlot(slot, dateLabel) {
    bookingState.slot = slot;
    document.getElementById('selected-slot-id').value = slot.ID;
    document.getElementById('people-count').max = slot.MaxPeople;
    document.getElementById('booking-form-container').classList.remove('hidden');

    document.getElementById('weather-info').innerHTML =
        '<h3 class="font-semibold mb-2">Детали бронирования</h3>' +
        '<p>Инструктор: <strong>' + bookingState.instructor.Name + '</strong></p>' +
        '<p>Прогулка: <strong>' + bookingState.walkType.Name + '</strong></p>' +
        '<p>Дата: <strong>' + dateLabel + '</strong></p>' +
        '<p>Время: <strong>' + slot.StartTime.substring(0, 5) + ' - ' + slot.EndTime.substring(0, 5) + '</strong></p>' +
        '<p>Цена: <strong class="text-blue-600">' + slot.Price + ' ₽</strong></p>' +
        '<p>Максимум: <strong>' + slot.MaxPeople + ' чел.</strong></p>';

    document.getElementById('booking-form-container').scrollIntoView({ behavior: 'smooth' });
}

function submitBookingForm(e) {
    e.preventDefault();
    var form = e.target;
    var formData = new FormData(form);
    var data = {
        slot_id: parseInt(formData.get('slot_id')),
        client_name: formData.get('client_name'),
        client_phone: formData.get('client_phone'),
        client_email: formData.get('client_email'),
        people_count: parseInt(formData.get('people_count'))
    };

    var btn = form.querySelector('button[type="submit"]');
    btn.disabled = true;
    btn.textContent = 'Отправка...';

    fetch('/booking', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    })
    .then(function(res) {
        if (!res.ok) {
            return res.text().then(function(err) { throw new Error(err || 'Ошибка при бронировании'); });
        }
        return res.json();
    })
    .then(function(result) {
        document.getElementById('booking-result').innerHTML = '<div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">' +
            '<p class="font-semibold">Бронирование отправлено!</p>' +
            '<p class="text-sm"><strong>Номер бронирования:</strong> #' + result.ID + '</p>' +
            '<p class="text-sm"><strong>Маршрут:</strong> ' + bookingState.walkType.Name + '</p>' +
            '<p class="text-sm mt-2"><strong>Статус:</strong> Ожидает подтверждения администратором в течение ' + result.hold_minutes + ' минут.</p>' +
            '</div>';
        form.reset();
        selectWalkType(bookingState.walkType);
    })
    .catch(function(err) {
        document.getElementById('booking-result').innerHTML = '<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">' +
            '<p class="font-semibold">Ошибка при создании бронирования</p>' +
            '<p class="text-sm">' + (err.message || 'Пожалуйста, попробуйте еще раз') + '</p></div>';
    })
    .finally(function() {
        btn.disabled = false;
        btn.textContent = 'Забронировать';
    });
}
