"""YAML Formatter
 
This script helps to format YAML files and can be used to perform simple ops like removing keys.
"""

import argparse
from yaml import safe_load, dump

def remove_keys(data, keys_to_remove):
    '''removes the specified keys from the yaml'''
    if isinstance(data, dict):
        for key in list(data.keys()):
            if key in keys_to_remove:
                del data[key]
            else:
                remove_keys(data[key], keys_to_remove)
    elif isinstance(data, list):
        for item in data:
            remove_keys(item, keys_to_remove)


def main(input_file, output_file, keys_to_remove=[]):
    '''formats input yaml file and saves output. Optionally removes keys from yaml file.'''
    print(f"Reading file {input_file}...")
    with open(input_file, 'r') as f_in:
        data = safe_load(f_in)

    if keys_to_remove:
        print(f"Removing keys {keys_to_remove} from {input_file}...")
        remove_keys(data, set(keys_to_remove))

    with open(output_file, 'w') as f_out:
        dump(data, f_out)
    print(f"Saved output to {output_file}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Script for formatting YAML files')
    parser.add_argument('-i', '--input', type=str, required=True, help="input yaml file")
    parser.add_argument('-o', '--output', type=str, default="output.yaml", help="output yaml file")
    parser.add_argument('-r', '--remove-keys', nargs='+', type=str, help='key(s) to remove, e.g. x-stoplight')

    args = parser.parse_args()

    main(input_file=args.input, output_file=args.output, keys_to_remove=args.remove_keys)


