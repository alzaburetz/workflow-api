#!/bin/bash

records=$1

phone=$((100000000000+$1))


generate_post_data() 
{
    phone=$((100000000000+$1))
    cat <<EOF
    {
        "name":"Test",
        "phone":"$phone",
        "email":"@$1",
        "password":"testtest"
    }
EOF
}

while [ $records -gt 1 ]
do
    echo "Sending request"
    curl -X POST "localhost:3000/api/user/register" -d "$(generate_post_data $records)"
    records=$(($records-1))
done