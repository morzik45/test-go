CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);
CREATE TABLE authorizations (
    id SERIAL NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    is_authorized BOOLEAN NOT NULL DEFAULT TRUE,
    login_at TIMESTAMP NOT NULL DEFAULT NOW(),
    logout_at TIMESTAMP,
    session_token VARCHAR(50) NOT NULL
);
CREATE TABLE variants (
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);
CREATE TABLE tasks (
    id SERIAL NOT NULL PRIMARY KEY,
    variant_id INT REFERENCES variants (id) ON DELETE CASCADE NOT NULL,
    question VARCHAR(255) NOT NULL,
    correct INT NOT NULL,
    answers JSON NOT NULL
);
CREATE TABLE tests (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    variant_id INT REFERENCES variants (id) ON DELETE CASCADE NOT NULL,
    start_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE TABLE user_answers (
    id SERIAL NOT NULL PRIMARY KEY,
    test_id INT REFERENCES tests (id) ON DELETE CASCADE NOT NULL,
    task_id INT REFERENCES tasks (id) ON DELETE CASCADE NOT NULL,
    answer INT NOT NULL
);
CREATE TABLE results (
    id SERIAL NOT NULL PRIMARY KEY,
    test_id INT REFERENCES tests (id) ON DELETE CASCADE NOT NULL,
    percent INT
);
INSERT INTO variants (name)
VALUES ('Математика');
INSERT INTO tasks (variant_id, question, correct, answers)
VALUES (
        (
            SELECT id
            from variants
            WHERE name = 'Математика'
        ),
        'Сколько будет 2+2?',
        2,
        '{"1":"3","2":"4","3":"6","4":"8"}'
    );
INSERT INTO tasks (variant_id, question, correct, answers)
VALUES (
        (
            SELECT id
            from variants
            WHERE name = 'Математика'
        ),
        'Фигура с 3 углами называется?',
        4,
        '{"1":"квадрат","2":"ромб","3":"катет","4":"треугольник"}'
    );
INSERT INTO tasks (variant_id, question, correct, answers)
VALUES (
        (
            SELECT id
            from variants
            WHERE name = 'Математика'
        ),
        'Сколько будет 3-8?',
        3,
        '{"1":"-1","2":"0","3":"-5","4":"-8"}'
    );
INSERT INTO variants (name)
VALUES ('Русский язык');
INSERT INTO tasks (variant_id, question, correct, answers)
VALUES (
        (
            SELECT id
            from variants
            WHERE name = 'Русский язык'
        ),
        'ЖИ ШИ пиши с?',
        2,
        '{"1":"ы","2":"и","3":"ё","4":"="}'
    );
INSERT INTO tasks (variant_id, question, correct, answers)
VALUES (
        (
            SELECT id
            from variants
            WHERE name = 'Русский язык'
        ),
        'Сколько букв Н в слове "деревя_ый"?',
        4,
        '{"1":"1","2":"0","3":"3","4":"2"}'
    );
INSERT INTO tasks (variant_id, question, correct, answers)
VALUES (
        (
            SELECT id
            from variants
            WHERE name = 'Русский язык'
        ),
        'перед "но" и "а" ставится?',
        3,
        '{"1":".","2":":","3":",","4":"-"}'
    );
INSERT INTO users (username, password_hash)
VALUES (
        'test',
        '686a7172686a7177313234363137616a6668616a7340bd001563085fc35165329ea1ff5c5ecbdbbeef'
    );