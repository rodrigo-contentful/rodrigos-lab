# Script to copy roles between two spaces
import contentful_management
import sys

#print(f"Arguments count: {len(sys.argv)}")
if len(sys.argv) < 2:
    print('missing arguments for origin and destination space')
    quit()
elif len(sys.argv) < 3:
    print('missing second argument for destination space')
    quit()

# get spaces id from arguments
orig_space=sys.argv[1]
dest_space=sys.argv[2]
cma_key='YOUR_CMA_KEY'

#conect to cma
client = contentful_management.Client(cma_key)

# get all roles on dest space and add them to array
my_input = [] 
dest_roles = client.roles(dest_space).all()
for dest_role in dest_roles:
    dest_role_name=getattr(dest_role, 'name', 'Name not found')
    print(f"Role found on destination '{dest_space}': {dest_role_name}")
    my_input.append(dest_role_name)

# on loop inserting roles, avoid already created roles 
roles = client.roles(orig_space).all()

for role in roles:
    role_name=getattr(role, 'name', 'Not a product')
    if role_name not in my_input:
        role_attributes = {
            'name': getattr(role, 'name', 'Name not found'),
            'description': getattr(role, 'description', 'description not found'),
            'permissions': getattr(role, 'permissions', 'permissions not found'),
            'policies': getattr(role, 'policies', 'policies not found'),
        }
        new_role = client.roles('fsnrcjskkd3y').create(role_attributes)
        print(f"Role crated on destination '{dest_space}': {role_name}")
